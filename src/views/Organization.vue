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
          <div class="branch-list">
            <div class="branch-card" v-for="branch in showData" :key="branch.branch">
              <div class="branch-header">
                <span class="branch-title">{{ branch.branch }}</span>
              </div>
              <table class="member-table">
                <thead>
                  <tr>
                    <th>姓名</th>
                    <th>职务</th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="member in branch.members" :key="member.name">
                    <td>{{ member.name }}</td>
                    <td>{{ member.duty }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, computed } from 'vue';
// mock 数据 - 按照组织结构图
const allData = [
  {
    branch: '2021级本科生党支部',
    members: [
      { name: '张三', duty: '党支部书记' },
      { name: '李四', duty: '党支部副书记' },
      { name: '王五', duty: '组织委员' },
      { name: '赵六', duty: '宣传委员' },
      { name: '孙七', duty: '纪检委员' },
      { name: '周八', duty: '统战委员' },
    ]
  },
  {
    branch: '2022级本科生党支部',
    members: [
      { name: '陈十', duty: '党支部书记' },
      { name: '吴十一', duty: '党支部副书记' },
      { name: '郑十二', duty: '组织委员' },
      { name: '王十三', duty: '宣传委员' },
      { name: '冯十四', duty: '纪检委员' },
      { name: '褚十五', duty: '统战委员' },
    ]
  }
];
const role = localStorage.getItem('role') || '党校校长';
const branch = localStorage.getItem('branch') || 'all';
const grade = localStorage.getItem('grade') || 'all';
const showData = computed(() => {
  if (role === '党校校长') {
    // 党校校长可以查看全部党支部
    return allData;
  } else if (role === '书记') {
    // 书记/副书记只能查看本年级的党支部
    return allData.map(item => ({
      ...item,
      grades: item.grades.filter(g => g.name === grade)
    })).filter(item => item.grades.length > 0);
  } else {
    // 默认显示全部
    return allData;
  }
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