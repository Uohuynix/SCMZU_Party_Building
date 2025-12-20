<template>
  <div class="page-layout">
    <aside class="sidebar">
      <router-link to="/dashboard">
        <img src="/src/pictures/dang.png" alt="党徽" class="sidebar-icon" />
        首页大屏
      </router-link>
      <router-link to="/development">
        <img src="/src/pictures/dang.png" alt="党徽" class="sidebar-icon" />
        党员档案
      </router-link>
      <router-link to="/training">
        <img src="/src/pictures/dang.png" alt="党徽" class="sidebar-icon" />
        培训信息
      </router-link>
      <router-link to="/organization">
        <img src="/src/pictures/dang.png" alt="党徽" class="sidebar-icon" />
        组织关系
      </router-link>
    </aside>
    <div class="main-content">
      <div class="page">
        <div class="card">
          <div class="card-header">
            <img class="logo" src="/src/pictures/dang.png" alt="党徽" />
            <h2>基层党支部组织架构</h2>
          </div>

          <!-- 筛选区域 -->
          <div class="filter-bar">
            <div class="filter-item">
              <label>党支部：</label>
              <select v-model="selectedBranch">
                <option value="all">全部</option>
                <option v-for="b in branchOptions" :key="b" :value="b">
                  {{ b }}
                </option>
              </select>
            </div>
            <div class="filter-item">
              <label>党小组：</label>
              <select v-model="selectedGroup">
                <option value="all">全部</option>
                <option v-for="g in groupOptions" :key="g" :value="g">
                  {{ g }}
                </option>
              </select>
            </div>
            <div class="filter-item">
              <label>发展阶段：</label>
              <select v-model="selectedType">
                <option value="all">全部</option>
                <option value="masses">群众</option>
                <option value="activist">入党积极分子</option>
                <option value="candidate">发展对象</option>
                <option value="probationary">预备党员</option>
                <option value="formal">正式党员</option>
              </select>
            </div>
          </div>

          <div v-if="loading" class="loading-state">加载组织数据中...</div>
          <div v-else-if="error" class="error-state">{{ error }}</div>

          <!-- 按支部/小组分组展示学员 -->
          <div v-else class="branch-list">
            <div
              class="branch-card"
              v-for="branch in groupedData"
              :key="branch.branchName"
            >
              <div class="branch-header">
                <span class="branch-title">{{ branch.branchName }}</span>
                <div class="branch-info">
                  <span class="tag tag-gray">小组数：{{ branch.groups.length }}</span>
                  <span class="tag tag-gray">总人数：{{ branch.totalMembers }}</span>
                </div>
              </div>

              <div class="group-list">
                <div
                  class="group-card"
                  v-for="group in branch.groups"
                  :key="group.groupName"
                >
                  <div class="group-title">
                    {{ group.groupName || '未分配党小组' }}
                    <span class="group-leader">
                      共 {{ group.members.length }} 人
                    </span>
                  </div>
                  <div class="member-list">
                    <div
                      class="member"
                      v-for="member in group.members"
                      :key="member.id"
                    >
                      <span class="avatar"></span>
                      <span>{{ member.name }}</span>
                      <span class="member-tag">{{ formatType(member.type) }}</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-if="!groupedData.length" class="empty-row">
              当前筛选条件下暂无人员数据
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';

const API_BASE_URL = '/api/v1';
const token = localStorage.getItem('token') || '';
const userRole = localStorage.getItem('role') || 'student';
const currentBranch = localStorage.getItem('branch') || 'all';
const currentGroup = localStorage.getItem('groupName') || 'all';
const isAdmin = userRole === 'admin';

const students = ref([]);
const loading = ref(false);
const error = ref('');

// 筛选条件
// 管理员默认查看“全部”支部和小组，其它角色默认定位到自己所在支部/小组
const selectedBranch = ref(
  isAdmin ? 'all' : (currentBranch === 'all' ? 'all' : currentBranch)
);
const selectedGroup = ref(
  isAdmin ? 'all' : (currentGroup === 'all' ? 'all' : currentGroup)
);
const selectedType = ref('all');

// 获取学员数据（后端已经按角色做了范围控制）
const fetchStudents = async () => {
  loading.value = true;
  error.value = '';
  try {
    const resp = await axios.get(`${API_BASE_URL}/students`, {
      params: { page: 1, pageSize: 1000 },
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
    const items = Array.isArray(resp.data)
      ? resp.data
      : resp.data.data || [];
    students.value = items;
  } catch (e) {
    console.error('获取组织数据失败', e);
    error.value = '获取组织数据失败，请稍后重试';
    students.value = [];
  } finally {
    loading.value = false;
  }
};

// 党支部 / 小组选项
const branchOptions = computed(() => {
  const set = new Set();
  students.value.forEach((s) => {
    if (s.branch) set.add(s.branch);
  });
  return Array.from(set);
});

const groupOptions = computed(() => {
  const set = new Set();
  students.value.forEach((s) => {
    if (
      (!selectedBranch.value || selectedBranch.value === 'all' || s.branch === selectedBranch.value) &&
      s.group_name
    ) {
      set.add(s.group_name);
    }
  });
  return Array.from(set);
});

// 将学员按支部 / 小组分组
const groupedData = computed(() => {
  const map = new Map();

  students.value.forEach((s) => {
    // 过滤条件
    if (selectedBranch.value !== 'all' && s.branch !== selectedBranch.value) {
      return;
    }
    if (selectedGroup.value !== 'all' && s.group_name !== selectedGroup.value) {
      return;
    }
    if (selectedType.value !== 'all' && s.type !== selectedType.value) {
      return;
    }

    const branchName = s.branch || '未分配党支部';
    const groupName = s.group_name || '未分配党小组';

    if (!map.has(branchName)) {
      map.set(branchName, new Map());
    }
    const groupMap = map.get(branchName);
    if (!groupMap.has(groupName)) {
      groupMap.set(groupName, []);
    }
    groupMap.get(groupName).push({
      id: s.id,
      name: s.name,
      type: s.type
    });
  });

  // 转换为模板可用结构
  const result = [];
  map.forEach((groupMap, branchName) => {
    const groups = [];
    let total = 0;
    groupMap.forEach((members, groupName) => {
      groups.push({
        groupName,
        members
      });
      total += members.length;
    });
    result.push({
      branchName,
      groups,
      totalMembers: total
    });
  });

  // 管理员可以看到所有支部，教师 / 学生默认只看到本支部（后端也有限制）
  return result.sort((a, b) => a.branchName.localeCompare(b.branchName, 'zh-Hans-CN'));
});

// 发展阶段中文映射
const formatType = (type) => {
  const map = {
    masses: '群众',
    activist: '入党积极分子',
    candidate: '发展对象',
    probationary: '预备党员',
    formal: '正式党员'
  };
  return map[type] || '其他';
};

onMounted(() => {
  fetchStudents();
});
</script>

<style scoped>
.page-layout {
  display: flex;
  min-height: 100vh;
  background: #f7f7fa;
}
.sidebar {
  width: 220px;
  background: linear-gradient(180deg, #fff 0%, #f8f9fa 100%);
  border-right: 1px solid #e0e0e0;
  padding: 30px 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  position: fixed;
  top: 0;
  left: 0;
  height: 100%;
  z-index: 10;
  box-shadow: 2px 0 10px rgba(0,0,0,0.05);
}
.sidebar-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1px solid #b71c1c;
  background: #fff;
}
.sidebar a {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12px;
  padding: 16px 20px;
  color: #555;
  text-decoration: none;
  font-size: 15px;
  font-weight: 500;
  border-radius: 12px;
  margin: 0 15px;
  transition: all 0.3s ease;
  width: 180px;
  position: relative;
  border: 1px solid transparent;
}
.sidebar a:hover {
  background: linear-gradient(135deg, #fff5f5 0%, #ffeaea 100%);
  color: #b71c1c;
  border-color: #ffcdd2;
  transform: translateX(5px);
  box-shadow: 0 4px 12px rgba(183,28,28,0.15);
}
.sidebar a.router-link-active {
  background: linear-gradient(135deg, #b71c1c 0%, #d32f2f 100%);
  color: white;
  font-weight: 600;
  border-color: #b71c1c;
  box-shadow: 0 4px 15px rgba(183,28,28,0.3);
}
.sidebar a.router-link-active::before {
  content: '';
  position: absolute;
  left: -15px;
  top: 50%;
  transform: translateY(-50%);
  width: 4px;
  height: 20px;
  background: #b71c1c;
  border-radius: 2px;
}
.main-content {
  flex: 1;
  padding-left: 220px; /* Adjust based on sidebar width */
  width: calc(100% - 220px); /* Adjust based on sidebar width */
}
.page {
  min-height: 100vh;
  background: #f7f7fa;
  padding: 36px 0;
  display: flex;
  justify-content: center;
}
.card {
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
  border: 1px solid #ececec;
  padding: 28px 22px 18px 22px;
  min-width: 900px;
  max-width: 95vw;
}
.card-header {
  display: flex;
  align-items: center;
  gap: 14px;
  margin-bottom: 12px;
}
.filter-bar {
  display: flex;
  gap: 16px;
  align-items: center;
  margin-bottom: 16px;
}
.filter-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
}
.filter-item label {
  color: #555;
}
.filter-item select {
  padding: 4px 8px;
  border-radius: 4px;
  border: 1px solid #ccc;
}
.logo {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  border: 1.5px solid #e53935;
  background: #fff;
}
h2 {
  color: #b71c1c;
  font-size: 1.3rem;
  font-weight: bold;
  letter-spacing: 1px;
}
.branch-list {
  display: flex;
  flex-wrap: wrap;
  gap: 24px;
  justify-content: flex-start;
}
.branch-card {
  background: #f9f9fb;
  border-radius: 10px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
  border: 1px solid #e0e0e0;
  padding: 18px 20px;
  min-width: 340px;
  flex: 1 1 350px;
}
.branch-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ececec;
  padding-bottom: 8px;
  margin-bottom: 12px;
}
.branch-title {
  font-size: 16px;
  font-weight: bold;
  color: #b71c1c;
  letter-spacing: 1px;
}
.branch-info {
  display: flex;
  gap: 6px;
}
.tag {
  display: inline-block;
  padding: 1px 8px;
  border-radius: 10px;
  font-size: 12px;
  font-weight: 500;
  margin-right: 2px;
  background: #f5f5f5;
  color: #b71c1c;
  border: 1px solid #e0e0e0;
}
.tag-red {
  background: #ffeaea;
  color: #b71c1c;
  border: 1px solid #e57373;
}
.tag-gold {
  background: #fff8e1;
  color: #b71c1c;
  border: 1px solid #ffe082;
}
.tag-gray {
  background: #f0f0f0;
  color: #616161;
  border: 1px solid #bdbdbd;
}
.grade-list {
  margin-left: 8px;
}
.grade-card {
  background: #fff;
  border-radius: 8px;
  margin-bottom: 12px;
  padding: 10px 14px;
  border-left: 3px solid #b71c1c;
  box-shadow: 0 1px 4px rgba(0,0,0,0.04);
}
.grade-title {
  font-size: 14px;
  font-weight: bold;
  color: #b71c1c;
  margin-bottom: 6px;
}
.group-list {
  margin-left: 10px;
}
.group-card {
  background: #f9f9fb;
  border-radius: 6px;
  margin-bottom: 8px;
  padding: 8px 12px;
  box-shadow: 0 1px 2px rgba(0,0,0,0.03);
  border-left: 2px solid #e0e0e0;
}
.group-title {
  font-size: 13px;
  font-weight: 500;
  color: #b71c1c;
  margin-bottom: 4px;
}
.group-leader {
  color: #b71c1c;
  font-size: 12px;
  margin-left: 6px;
}
.member-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 4px;
}
.member {
  background: #fff;
  color: #b71c1c;
  border-radius: 12px;
  padding: 2px 12px 2px 6px;
  font-size: 13px;
  display: flex;
  align-items: center;
  box-shadow: 0 1px 2px rgba(0,0,0,0.03);
  font-weight: 500;
  border: 1px solid #ececec;
}
.avatar {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  background: #f5f5f5;
  margin-right: 6px;
  border: 1px solid #e0e0e0;
  display: inline-block;
}
.member-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 10px;
}
.member-table th,
.member-table td {
  padding: 8px 12px;
  text-align: left;
  border-bottom: 1px solid #ececec;
}
.member-table th {
  background-color: #f5f5f5;
  font-weight: bold;
  color: #333;
}
.member-table tr:last-child td {
  border-bottom: none;
}
</style> 