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
        <div class="container">
          <aside class="member-list-panel">
            <div class="search-bar">
              <input v-model="search" placeholder="搜索姓名/学号/班级..." />
            </div>
            <div class="import-export-bar">
              <input type="file" ref="fileInput" style="display:none" @change="handleImport" accept=".xlsx,.xls" />
              <button class="btn" @click="triggerImport">导入Excel</button>
              <button class="btn btn-primary" @click="exportExcel">导出Excel</button>
            </div>
            <div v-if="loading && members.length === 0" class="member-list-loading">
              加载党员列表中...
            </div>
            <div v-else-if="error && members.length === 0" class="member-list-loading error">
              {{ error }}
            </div>
            <ul v-else class="member-list">
              <li v-for="member in filteredMembers" :key="member.id" :class="{active: member.id === selectedId}" @click="selectMember(member.id)">
                <span class="name">{{ member.name }}</span>
                <span class="id">{{ member.sid }}</span>
                <span class="class">{{ member.class }}</span>
              </li>
              <li v-if="filteredMembers.length === 0" class="empty-list">
                未找到匹配的党员
              </li>
            </ul>
          </aside>
          <main class="archive-main">
            <div v-if="selectedMember" class="member-card">
              <div class="card-header">
                <div class="header-left">
                  <div class="header-icon">!</div>
                  <h2>党员信息卡</h2>
                </div>
                <div class="header-close">×</div>
              </div>
              
              <div class="tab-navigation">
                <div class="tab" :class="{ active: activeTab === 'basic' }" @click="activeTab = 'basic'">基本信息</div>
                <div class="tab" :class="{ active: activeTab === 'development' }" @click="activeTab = 'development'">党员发展</div>
              </div>
              
              <div class="card-content">
                <!-- 基本信息 -->
                <div v-show="activeTab === 'basic'" class="tab-content">
                  <div v-if="loading" class="loading-state">加载中...</div>
                  <div v-else-if="selectedMember" class="info-section">
                    <div class="info-left">
                      <div class="info-row">
                        <span class="info-label">姓名</span>
                        <span class="info-value highlight">{{ selectedMember.name }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">性别</span>
                        <span class="info-value">{{ selectedMember.gender }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">民族</span>
                        <span class="info-value">{{ selectedMember.nation }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">出生日期</span>
                        <span class="info-value">{{ selectedMember.birth }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">学历</span>
                        <span class="info-value">{{ selectedMember.edu }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">联系方式</span>
                        <span class="info-value">{{ selectedMember.phone }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">身份证号</span>
                        <span class="info-value">{{ selectedMember.idcard }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">对象类型</span>
                        <span class="info-value">{{ selectedMember.memberType }}</span>
                      </div>
                    </div>
                    
                    <div class="info-right">
                      <div class="info-row">
                        <span class="info-label">学号/教工号</span>
                        <span class="info-value">{{ selectedMember.sid }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">班级/单位</span>
                        <span class="info-value">{{ selectedMember.class }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">年级</span>
                        <span class="info-value">{{ selectedMember.grade }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">政治面貌</span>
                        <span class="info-value">{{ selectedMember.memberType }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">入党时间</span>
                        <span class="info-value">{{ selectedMember.applyDate }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">转正时间</span>
                        <span class="info-value">{{ selectedMember.formalDate }}</span>
                      </div>
                    </div>
                    
                    <div class="photo-section">
                      <div class="photo-placeholder">
                        <div class="photo-silhouette"></div>
                        <span class="photo-text">暂无照片</span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- 党员发展 -->
                <div v-show="activeTab === 'development'" class="tab-content">
                  <div v-if="loading" class="loading-state">加载中...</div>
                  <div v-else-if="selectedMember" class="development-section">
                    <div class="info-left">
                      <div class="info-row">
                        <span class="info-label">入党申请人时间</span>
                        <span class="info-value">{{ selectedMember.applyDate || (memberDevelopment?.applyDate ? new Date(memberDevelopment.applyDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">积极分子时间</span>
                        <span class="info-value">{{ selectedMember.positiveDate || (memberDevelopment?.positiveDate ? new Date(memberDevelopment.positiveDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">发展对象时间</span>
                        <span class="info-value">{{ selectedMember.devDate || (memberDevelopment?.devDate ? new Date(memberDevelopment.devDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">预备党员时间</span>
                        <span class="info-value">{{ selectedMember.preDate || (memberDevelopment?.preDate ? new Date(memberDevelopment.preDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                    </div>
                    
                    <div class="info-right">
                      <div class="info-row">
                        <span class="info-label">转正时间</span>
                        <span class="info-value">{{ selectedMember.formalDate || (memberDevelopment?.formalDate ? new Date(memberDevelopment.formalDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">党组织关系转接时间</span>
                        <span class="info-value">{{ selectedMember.orgTransDate || (memberDevelopment?.orgTransDate ? new Date(memberDevelopment.orgTransDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">党组织关系介绍信时间</span>
                        <span class="info-value">{{ selectedMember.orgIntroDate || (memberDevelopment?.orgIntroDate ? new Date(memberDevelopment.orgIntroDate).toISOString().split('T')[0] : '') }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">发展状态</span>
                        <span class="info-value highlight">{{ selectedMember.memberType === '正式党员' ? '已转正' : '未转正' }}</span>
                      </div>
                    </div>
                  </div>
                </div>


                

              </div>
            </div>
            <div v-else class="empty-tip">请选择左侧党员查看详细档案</div>
          </main>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import * as XLSX from 'xlsx';
import axios from 'axios';

const fileInput = ref(null);
// 数据状态
const members = ref([]);
const search = ref('');
const selectedId = ref(null);
const selectedMember = ref(null);
const memberDevelopment = ref(null);
const loading = ref(false);
const error = ref(null);
const activeTab = ref('basic');

// API基础URL
const API_BASE_URL = 'http://127.0.0.1:4523/m1/7299957-7028709-default/api/v1';

// 过滤后的党员列表
const filteredMembers = computed(() => {
  if (!search.value) return members.value;
  return members.value.filter(m => 
    m.name?.includes(search.value) || 
    m.student_no?.includes(search.value) || 
    m.major_class?.includes(search.value)
  );
});

// 数据映射函数
const mapStudentToMember = (student) => ({
  id: student.id,
  name: student.name,
  sid: student.student_no,
  class: student.major_class,
  gender: student.gender === 'male' ? '男' : '女',
  nation: student.ethnicity,
  birth: student.birth_date ? new Date(student.birth_date).toISOString().split('T')[0] : '',
  edu: student.education,
  phone: student.phone,
  idcard: student.id_card,
  grade: student.admission_date ? new Date(student.admission_date).getFullYear() + '级' : '',
  memberType: mapPartyType(student.type),
  applyDate: student.admission_date ? new Date(student.admission_date).toISOString().split('T')[0] : '',
  positiveDate: '',
  devDate: '',
  preDate: '',
  formalDate: student.conversion_date ? new Date(student.conversion_date).toISOString().split('T')[0] : '',
  orgTransDate: '',
  orgIntroDate: ''
});

// 党类型映射
const mapPartyType = (type) => {
  const typeMap = {
    'masses': '群众',
    'party_member': '正式党员',
    'probationary_party_member': '预备党员',
    'activist': '入党积极分子'
  };
  return typeMap[type] || type;
};

// 获取党员列表
const fetchMembers = async () => {
  loading.value = true;
  error.value = null;
  try {
    const response = await axios.get(`${API_BASE_URL}/students`);
    // 根据API返回的数据结构可能需要调整
    if (Array.isArray(response.data)) {
      members.value = response.data.map(student => mapStudentToMember(student));
    } else if (response.data.data && Array.isArray(response.data.data)) {
      members.value = response.data.data.map(student => mapStudentToMember(student));
    } else {
      members.value = [];
    }
    
    // 如果有数据且未选择任何党员，默认选择第一个
    if (members.value.length > 0 && selectedId.value === null) {
      selectMember(members.value[0].id);
    }
  } catch (err) {
    error.value = '获取党员列表失败';
    console.error('获取党员列表失败:', err);
    // 使用模拟数据作为备用
    members.value = [
      { id: 1, name: '张三', sid: '2021001', class: '软件工程1班', gender: '男', nation: '汉族', birth: '1999-01-01', edu: '本科', phone: '13800000000', idcard: '1234567890', grade: '2018级', memberType: '正式党员', applyDate: '2017-10-01', formalDate: '2019-09-01' },
      { id: 2, name: '李四', sid: '2021002', class: '软件工程2班', gender: '女', nation: '汉族', birth: '2000-02-02', edu: '硕士', phone: '13800000001', idcard: '1234567891', grade: '2019级', memberType: '预备党员', applyDate: '2018-10-01', formalDate: '2020-09-01' },
    ];
  } finally {
    loading.value = false;
  }
};

// 获取党员基本信息
const fetchMemberBasicInfo = async (id) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/students/${id}`);
    const studentData = response.data.data || response.data;
    const mappedData = mapStudentToMember(studentData);
    
    // 更新选中的党员信息
    selectedMember.value = mappedData;
    
    // 同时更新members数组中的对应项
    const index = members.value.findIndex(m => m.id === id);
    if (index !== -1) {
      members.value[index] = mappedData;
    }
  } catch (err) {
    console.error('获取党员基本信息失败:', err);
    // 如果API失败，保持当前信息不变
  }
};

// 获取党员发展信息
const fetchMemberDevelopment = async (id) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/students/${id}/development`);
    memberDevelopment.value = response.data.data || response.data;
  } catch (err) {
    console.error('获取党员发展信息失败:', err);
    memberDevelopment.value = null;
  }
};
// 选择党员
async function selectMember(id) {
  selectedId.value = id;
  loading.value = true;
  
  // 并行获取基本信息和发展信息
  try {
    await Promise.all([
      fetchMemberBasicInfo(id),
      fetchMemberDevelopment(id)
    ]);
  } catch (err) {
    console.error('加载党员信息失败:', err);
  } finally {
    loading.value = false;
  }
}

// 组件挂载时获取党员列表
onMounted(() => {
  fetchMembers();
});
function triggerImport() {
  fileInput.value && fileInput.value.click();
}
function handleImport(e) {
  const file = e.target.files[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = (evt) => {
    const data = new Uint8Array(evt.target.result);
    const workbook = XLSX.read(data, { type: 'array' });
    const sheet = workbook.Sheets[workbook.SheetNames[0]];
    const json = XLSX.utils.sheet_to_json(sheet);
    // 只导入入党申请人及后续信息
    const newMembers = json.map((row, idx) => ({
      id: Date.now() + idx,
      name: row['姓名'] || '',
      sid: row['学号/工号'] || '',
      class: row['班级/单位'] || '',
      gender: row['性别'] || '',
      nation: row['民族'] || '',
      birth: row['出生日期'] || '',
      edu: row['学历'] || '',
      phone: row['联系方式'] || '',
      idcard: row['身份证号'] || '',
      grade: row['年级'] || '',
      memberType: ['正式党员','预备党员','发展对象','入党积极分子','入党申请人'].includes(row['对象类型']) ? row['对象类型'] : '入党申请人',
      applyDate: row['入党申请人时间'] || '',
      positiveDate: row['积极分子时间'] || '',
      devDate: row['发展对象时间'] || '',
      preDate: row['预备党员时间'] || '',
      formalDate: row['转正时间'] || '',
      orgTransDate: row['党组织关系转接时间'] || '',
      orgIntroDate: row['党组织关系介绍信时间'] || '',
    }));
    members.value = [...members.value, ...newMembers];
    if (newMembers.length > 0) selectedId.value = newMembers[0].id;
  };
  reader.readAsArrayBuffer(file);
}
function exportExcel() {
  const data = members.value.map(m => ({
    姓名: m.name,
    '学号/工号': m.sid,
    性别: m.gender,
    民族: m.nation,
    出生日期: m.birth,
    学历: m.edu,
    联系方式: m.phone,
    身份证号: m.idcard,
    '班级/单位': m.class,
    年级: m.grade,
    对象类型: m.memberType,
    入党申请人时间: m.applyDate,
    积极分子时间: m.positiveDate,
    发展对象时间: m.devDate,
    预备党员时间: m.preDate,
    转正时间: m.formalDate,
    党组织关系转接时间: m.orgTransDate,
    党组织关系介绍信时间: m.orgIntroDate,
  }));
  const ws = XLSX.utils.json_to_sheet(data);
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, '党员档案');
  XLSX.writeFile(wb, '党员档案.xlsx');
}
</script>

<style scoped>
/* 加载状态样式 */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 300px;
  font-size: 16px;
  color: #666;
}

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
.container {
  display: flex;
  width: 1200px;
  background: #fff;
  border-radius: 14px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.07);
  border: 1px solid #ececec;
  overflow: hidden;
}
.member-list-panel {
  width: 260px;
  background: #f9f9fb;
  border-right: 1px solid #ececec;
  display: flex;
  flex-direction: column;
}
.search-bar {
  padding: 18px 16px 8px 16px;
  background: #fff;
}
.search-bar input {
  width: calc(100% - 24px);
  padding: 7px 12px;
  border-radius: 6px;
  border: 1px solid #e0e0e0;
  font-size: 14px;
  box-sizing: border-box;
}
/* 党员列表加载状态 */
  .member-list-loading {
    display: flex;
    justify-content: center;
    align-items: center;
    padding: 40px 20px;
    color: #666;
  }
  
  /* 空列表状态 */
  .empty-list {
    text-align: center;
    padding: 20px;
    color: #999;
    font-style: italic;
  }
.btn {
  background: #fff;
  color: #b71c1c;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 5px 14px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: background 0.2s, color 0.2s, border 0.2s;
}
.btn:hover {
  background: #fbe9e7;
  color: #b71c1c;
  border: 1px solid #b71c1c;
}
.btn-primary {
  background: #b71c1c;
  color: #fff;
  border: none;
}
.member-list {
  flex: 1;
  overflow-y: auto;
  padding: 0 0 12px 0;
  margin: 0;
  list-style: none;
}
.member-list li {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 18px;
  cursor: pointer;
  border-left: 3px solid transparent;
  transition: background 0.2s, border 0.2s;
}
.member-list li.active, .member-list li:hover {
  background: #fbe9e7;
  border-left: 3px solid #b71c1c;
}

.name {
  color: #b71c1c;
  font-weight: bold;
  font-size: 15px;
}
.id, .class {
  color: #616161;
  font-size: 13px;
  margin-left: 6px;
}
.archive-main {
  flex: 1;
  padding: 20px;
}

/* 党员信息卡片样式 */
.member-card {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  overflow: hidden;
  max-width: 900px;
  margin: 0 auto;
}

.card-header {
  background: #b71c1c;
  color: white;
  padding: 15px 20px;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 10px;
}

.header-icon {
  width: 24px;
  height: 24px;
  background: white;
  color: #b71c1c;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 14px;
}

.header-close {
  font-size: 20px;
  cursor: pointer;
  padding: 5px;
}

.tab-navigation {
  display: flex;
  background: #f5f5f5;
  border-bottom: 1px solid #e0e0e0;
}

.tab {
  padding: 12px 20px;
  cursor: pointer;
  border-bottom: 3px solid transparent;
  transition: all 0.2s;
  background: #f5f5f5;
  color: #333;
}

.tab:hover {
  background: #e0e0e0;
  color: #b71c1c;
}

.tab.active {
  background: white;
  border-bottom-color: #b71c1c;
  color: #b71c1c;
  font-weight: 500;
}

.tab-content {
  min-height: 300px;
}

.development-section {
  display: flex;
  gap: 60px;
  position: relative;
  align-items: flex-start;
}

.development-section .info-left, .development-section .info-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 18px;
  min-width: 0;
}

.development-section .info-row {
  display: flex;
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid #f5f5f5;
  min-height: 24px;
}

.development-section .info-label {
  width: 140px;
  color: #666;
  font-size: 14px;
  font-weight: 500;
  flex-shrink: 0;
}

.development-section .info-value {
  flex: 1;
  color: #333;
  font-size: 14px;
  padding-left: 10px;
}

/* 已移除奖惩信息和党内职务相关样式 */

.card-content {
  padding: 30px;
}

.info-section {
  display: flex;
  gap: 40px;
  position: relative;
}

.info-left, .info-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.info-row {
  display: flex;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-label {
  width: 100px;
  color: #666;
  font-size: 14px;
  font-weight: 500;
}

.info-value {
  flex: 1;
  color: #333;
  font-size: 14px;
}

.info-value.highlight {
  color: #b71c1c;
  font-weight: bold;
}

.photo-section {
  position: absolute;
  right: 0;
  top: 0;
  width: 120px;
  height: 160px;
}

.photo-placeholder {
  width: 100%;
  height: 100%;
  background: #f5f5f5;
  border: 2px dashed #ccc;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.photo-silhouette {
  width: 60px;
  height: 60px;
  background: #ddd;
  border-radius: 50%;
  position: relative;
}

.photo-silhouette::before {
  content: '';
  position: absolute;
  top: 15px;
  left: 50%;
  transform: translateX(-50%);
  width: 30px;
  height: 30px;
  background: #bbb;
  border-radius: 50%;
}

.photo-silhouette::after {
  content: '';
  position: absolute;
  bottom: 5px;
  left: 50%;
  transform: translateX(-50%);
  width: 40px;
  height: 20px;
  background: #bbb;
  border-radius: 20px 20px 0 0;
}

.photo-text {
  font-size: 12px;
  color: #999;
}

.action-section {
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
}

.edit-btn {
  background: #b71c1c;
  color: white;
  border: none;
  padding: 10px 24px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(183,28,28,0.2);
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.edit-btn::before {
  content: "✏️";
  font-size: 16px;
}

.edit-btn:hover {
  background: #d32f2f;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(183,28,28,0.3);
}
.empty-tip {
  color: #bdbdbd;
  font-size: 16px;
  text-align: center;
  margin-top: 80px;
}
</style>