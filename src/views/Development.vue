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
            <!-- 导入/导出仅管理员和教师可见 -->
            <div class="import-export-bar" v-if="userRole !== 'student'">
              <input type="file" ref="fileInput" style="display:none" @change="handleImport" accept=".xlsx,.xls" />
              <button class="btn" @click="triggerImport">导入Excel</button>
              <button class="btn btn-primary" @click="exportExcel">导出Excel</button>
            </div>
            <!-- 学生：如果还没有自己的党员档案，提供创建按钮 -->
            <div
              class="import-export-bar"
              v-else-if="userRole === 'student' && !hasOwnRecord"
            >
              <button class="btn btn-primary" @click="createSelfRecord">
                创建我的党员档案
              </button>
            </div>
            <div v-if="loading && members.length === 0" class="member-list-loading">
              加载党员列表中...
            </div>
            <div v-else-if="error && members.length === 0" class="member-list-loading error">
              {{ error }}
            </div>
            <ul v-else class="member-list">
              <li
                v-for="member in filteredMembers"
                :key="member.id"
                :class="{active: member.id === selectedId}"
                @click="selectMember(member.id)"
              >
                <span class="name">{{ member.name }}</span>
                <span class="id">{{ member.sid }}</span>
                <span class="class">{{ member.class }}</span>
              </li>
              <li v-if="filteredMembers.length === 0" class="empty-list">
                <span v-if="userRole === 'student'">
                  暂无党员档案，点击上方“创建我的党员档案”按钮进行录入
                </span>
                <span v-else>
                  未找到匹配的党员
                </span>
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
                        <span v-if="!editMode" class="info-value highlight">{{ selectedMember.name }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.name"
                          type="text"
                        />
                      </div>
                      <div class="info-row">
                        <span class="info-label">性别</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.gender }}</span>
                        <select v-else class="info-value" v-model="editForm.gender">
                          <option value="男">男</option>
                          <option value="女">女</option>
                        </select>
                      </div>
                      <div class="info-row">
                        <span class="info-label">民族</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.nation }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.nation"
                          type="text"
                        />
                      </div>
                      <div class="info-row">
                        <span class="info-label">出生日期</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.birth }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.birth"
                          type="date"
                        />
                      </div>
                      <div class="info-row">
                        <span class="info-label">学历</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.edu }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.edu"
                          type="text"
                        />
                      </div>
                      <div class="info-row">
                        <span class="info-label">联系方式</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.phone }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.phone"
                          type="text"
                        />
                      </div>
                      <div class="info-row">
                        <span class="info-label">身份证号</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.idcard }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.idcard"
                          type="text"
                        />
                      </div>
                      <div class="info-row">
                        <span class="info-label">对象类型</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.memberType }}</span>
                        <select v-else class="info-value" v-model="editForm.memberType">
                          <option value="入党申请人">入党申请人</option>
                          <option value="入党积极分子">入党积极分子</option>
                          <option value="发展对象">发展对象</option>
                          <option value="预备党员">预备党员</option>
                          <option value="正式党员">正式党员</option>
                        </select>
                      </div>
                    </div>
                    
                    <div class="info-right">
                      <div class="info-row">
                        <span class="info-label">学号/教工号</span>
                        <span class="info-value">{{ selectedMember.sid }}</span>
                      </div>
                      <div class="info-row">
                        <span class="info-label">班级/单位</span>
                        <span v-if="!editMode" class="info-value">{{ selectedMember.class }}</span>
                        <input
                          v-else
                          class="info-value"
                          v-model="editForm.class"
                          type="text"
                        />
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
                      <div v-if="currentPhotoUrl" class="photo-wrapper">
                        <img :src="currentPhotoUrl" alt="党员照片" class="photo-image" />
                      </div>
                      <div v-else class="photo-placeholder">
                        <div class="photo-silhouette"></div>
                        <span class="photo-text">暂无照片</span>
                      </div>
                      <div class="photo-actions" v-if="canEditCurrent">
                        <input
                          ref="photoInput"
                          type="file"
                          accept="image/*"
                          style="display:none"
                          @change="handlePhotoChange"
                        />
                        <button class="btn" type="button" @click="triggerPhotoSelect">
                          {{ currentPhotoUrl ? '更换照片' : '上传照片' }}
                        </button>
                      </div>
                    </div>
                  </div>

                  <!-- 操作区域：根据权限显示编辑按钮 -->
                  <div class="action-section" v-if="canEditCurrent">
                    <button
                      v-if="!editMode"
                      class="edit-btn"
                      type="button"
                      @click="startEdit"
                    >
                      编辑个人信息
                    </button>
                    <div v-else style="display:flex; gap:12px; margin-top:12px;">
                      <button
                        class="edit-btn"
                        type="button"
                        :disabled="saving"
                        @click="saveEdit"
                      >
                        {{ saving ? '保存中...' : '保存' }}
                      </button>
                      <button
                        class="btn"
                        type="button"
                        :disabled="saving"
                        @click="cancelEdit"
                      >
                        取消
                      </button>
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
const photoInput = ref(null);
const userRole = localStorage.getItem('role') || 'student';
const currentUsername = localStorage.getItem('username') || '';
// 数据状态
const members = ref([]);
const search = ref('');
const selectedId = ref(null);
const selectedMember = ref(null);
const memberDevelopment = ref(null);
const loading = ref(false);
const error = ref(null);
const activeTab = ref('basic');

// 编辑状态
const editMode = ref(false);
const saving = ref(false);
const editForm = ref({
  name: '',
  gender: '',
  nation: '',
  birth: '',
  edu: '',
  phone: '',
  idcard: '',
  class: '',
  memberType: ''
});

// 照片本地存储（仅前端预览，使用学号作为 key）
const photoStore = ref({});
const currentPhotoUrl = computed(() => {
  if (!selectedMember.value) return '';
  // 优先使用后端的 photo_url，如果没有则退回本地缓存（兼容旧数据）
  if (selectedMember.value.photoUrl) {
    return selectedMember.value.photoUrl;
  }
  const sid = selectedMember.value.sid;
  return photoStore.value[sid] || '';
});

// API基础URL（通过 Vite 代理到本地后端）
const API_BASE_URL = '/api/v1';
const token = localStorage.getItem('token') || '';

// 学生是否已经有自己的党员档案
const hasOwnRecord = computed(() =>
  userRole === 'student'
    ? members.value.some(m => m.sid === currentUsername)
    : false
);

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
  orgIntroDate: '',
  photoUrl: student.photo_url || ''
});

// 党类型映射（后端 -> 前端展示）
const mapPartyType = (type) => {
  const typeMap = {
    masses: '入党申请人',
    formal: '正式党员',
    probationary: '预备党员',
    candidate: '发展对象',
    activist: '入党积极分子'
  };
  return typeMap[type] || type;
};

// 党类型映射（前端展示 -> 后端存储）
const mapMemberTypeToPartyType = (memberType) => {
  const map = {
    '群众': 'masses',          // 兼容旧文案
    '入党申请人': 'masses',
    '入党积极分子': 'activist',
    '发展对象': 'candidate',
    '预备党员': 'probationary',
    '正式党员': 'formal'
  };
  return map[memberType] || 'masses';
};

// 是否允许当前登录用户编辑选中人员信息
const canEditCurrent = computed(() => {
  if (!selectedMember.value) return false;
  if (userRole === 'admin') return true;
  if (userRole === 'student') {
    // 学生只能编辑自己的信息（用户名约定为学号）
    return selectedMember.value.sid === currentUsername;
  }
  if (userRole === 'teacher') {
    // 教师可编辑当前列表中可见的成员，后端会再次做支部/小组校验
    return true;
  }
  return false;
});

// 获取党员列表
const fetchMembers = async () => {
  loading.value = true;
  error.value = null;
  try {
    const response = await axios.get(`${API_BASE_URL}/students`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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

// 学生创建自己的党员档案（最小信息，后续可在右侧编辑补充）
const createSelfRecord = async () => {
  if (userRole !== 'student') return;
  if (!currentUsername) {
    alert('当前登录账号缺少学号信息，无法创建档案');
    return;
  }
  if (hasOwnRecord.value) {
    alert('您已经有党员档案了');
    return;
  }

  const name = window.prompt('请输入您的姓名，用于创建党员档案：');
  if (!name) {
    return;
  }

  const payload = {
    student_no: currentUsername,
    name,
    // 性别后端是 ENUM('male','female')，默认 male 避免空值报错
    gender: 'male',
    ethnicity: '',
    birth_date: null,
    education: '',
    phone: '',
    id_card: '',
    major_class: '',
    type: mapMemberTypeToPartyType('入党申请人')
  };

  try {
    await axios.post(`${API_BASE_URL}/students`, payload, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
    // 重新拉取列表并自动选中自己的档案
    await fetchMembers();
    const mine = members.value.find(m => m.sid === currentUsername);
    if (mine) {
      await selectMember(mine.id);
    }
    alert('党员档案已创建，请在右侧补充完善个人信息');
  } catch (err) {
    console.error('创建本人党员档案失败:', err);
    const backendMsg =
      err?.response?.data?.error ||
      err?.message ||
      '创建党员档案失败，请稍后重试或联系管理员';
    alert(`创建党员档案失败：${backendMsg}`);
  }
};

// 获取党员基本信息
const fetchMemberBasicInfo = async (id) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/students/${id}`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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
    const response = await axios.get(`${API_BASE_URL}/students/${id}/development`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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
  editMode.value = false;
  
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
  // 读取本地已保存的照片映射
  try {
    const raw = localStorage.getItem('memberPhotos') || '{}';
    photoStore.value = JSON.parse(raw);
  } catch (e) {
    photoStore.value = {};
  }
  fetchMembers();
});
function triggerImport() {
  fileInput.value && fileInput.value.click();
}

// 根据表格行构建后端学生对象
const buildStudentPayloadFromRow = (row) => {
  const genderText = row['性别'] || '';
  const gender =
    genderText === '男' ? 'male' :
    genderText === '女' ? 'female' : '';

  const birthRaw = row['出生日期'] || '';
  const birthDate = birthRaw ? new Date(birthRaw) : null;

  const memberType =
    ['正式党员','预备党员','发展对象','入党积极分子','入党申请人'].includes(row['对象类型'])
      ? row['对象类型']
      : '入党申请人';

  return {
    student_no: row['学号/工号'] || '',
    name: row['姓名'] || '',
    gender,
    ethnicity: row['民族'] || '',
    birth_date: birthDate ? birthDate.toISOString() : null,
    education: row['学历'] || '',
    phone: row['联系方式'] || '',
    id_card: row['身份证号'] || '',
    major_class: row['班级/单位'] || '',
    // 这里暂不强依赖年级和入学时间，可根据需要扩展
    type: mapMemberTypeToPartyType(memberType),
  };
};

// 按学号做“有则更新，无则新增”
const upsertStudentFromRow = async (row) => {
  const payload = buildStudentPayloadFromRow(row);
  if (!payload.student_no) {
    return;
  }

  // 先在当前已加载列表中按学号查找
  const existing = members.value.find(m => m.sid === payload.student_no);

  try {
    if (existing && existing.id) {
      // 已存在该学号，用 PUT 更新
      await axios.put(`${API_BASE_URL}/students/${existing.id}`, payload, {
        headers: {
          Authorization: token ? `Bearer ${token}` : ''
        }
      });
    } else {
      // 不存在该学号，用 POST 新建
      await axios.post(`${API_BASE_URL}/students`, payload, {
        headers: {
          Authorization: token ? `Bearer ${token}` : ''
        }
      });
    }
  } catch (err) {
    console.error('导入单行失败:', payload.student_no, err);
  }
};

async function handleImport(e) {
  const file = e.target.files[0];
  if (!file) return;
  const reader = new FileReader();
  reader.onload = async (evt) => {
    const data = new Uint8Array(evt.target.result);
    const workbook = XLSX.read(data, { type: 'array' });
    const sheet = workbook.Sheets[workbook.SheetNames[0]];
    const json = XLSX.utils.sheet_to_json(sheet);

    // 逐行与后端同步，学号作为主键避免重复创建
    for (const row of json) {
      await upsertStudentFromRow(row);
    }

    // 同步完成后重新从后端拉取列表，确保刷新后/切换页面数据仍然存在
    await fetchMembers();

    if (members.value.length > 0) {
      selectedId.value = members.value[0].id;
      selectedMember.value = members.value[0];
    }
  };
  reader.readAsArrayBuffer(file);
}

// 触发选择照片
function triggerPhotoSelect() {
  photoInput.value && photoInput.value.click();
}

// 处理照片选择：上传到后端并更新当前记录的 photoUrl
async function handlePhotoChange(e) {
  const file = e.target.files[0];
  if (!file || !selectedMember.value) return;

  const formData = new FormData();
  formData.append('photo', file);

  try {
    const resp = await axios.post(
      `${API_BASE_URL}/students/${selectedMember.value.id}/photo`,
      formData,
      {
        headers: {
          Authorization: token ? `Bearer ${token}` : ''
          // 不要显式设置 Content-Type，浏览器会自动带 multipart 边界
        }
      }
    );
    const url = resp.data.photo_url;
    if (url) {
      // 更新当前选中成员的照片地址
      selectedMember.value = {
        ...selectedMember.value,
        photoUrl: url
      };
      // 为了兼容旧本地缓存逻辑，仍保存一份
      const sid = selectedMember.value.sid;
      photoStore.value = {
        ...photoStore.value,
        [sid]: url
      };
      localStorage.setItem('memberPhotos', JSON.stringify(photoStore.value));
    }
    alert('照片上传成功');
  } catch (err) {
    console.error('照片上传失败:', err);
    alert('照片上传失败，请稍后重试');
  } finally {
    // 重置 input，避免同一个文件无法再次触发 change
    if (photoInput.value) {
      photoInput.value.value = '';
    }
  }
}

// 初始化编辑表单
function startEdit() {
  if (!selectedMember.value) return;
  editForm.value = {
    name: selectedMember.value.name || '',
    gender: selectedMember.value.gender || '',
    nation: selectedMember.value.nation || '',
    birth: selectedMember.value.birth || '',
    edu: selectedMember.value.edu || '',
    phone: selectedMember.value.phone || '',
    idcard: selectedMember.value.idcard || '',
    class: selectedMember.value.class || '',
    memberType: selectedMember.value.memberType || ''
  };
  editMode.value = true;
}

function cancelEdit() {
  editMode.value = false;
}

// 保存编辑后的个人信息
async function saveEdit() {
  if (!selectedMember.value) return;
  saving.value = true;
  try {
    const genderText = editForm.value.gender || '';
    const gender =
      genderText === '男' ? 'male' :
      genderText === '女' ? 'female' : '';

    const birthRaw = editForm.value.birth || '';
    const birthDate = birthRaw ? new Date(birthRaw) : null;

    const payload = {
      name: editForm.value.name || '',
      gender,
      ethnicity: editForm.value.nation || '',
      birth_date: birthDate ? birthDate.toISOString() : null,
      education: editForm.value.edu || '',
      phone: editForm.value.phone || '',
      id_card: editForm.value.idcard || '',
      major_class: editForm.value.class || '',
      type: mapMemberTypeToPartyType(editForm.value.memberType || '入党申请人')
    };

    await axios.put(
      `${API_BASE_URL}/students/${selectedMember.value.id}`,
      payload,
      {
        headers: {
          Authorization: token ? `Bearer ${token}` : ''
        }
      }
    );

    // 重新加载当前成员信息，确保界面数据最新
    await fetchMemberBasicInfo(selectedMember.value.id);
    editMode.value = false;
    alert('个人信息已保存');
  } catch (err) {
    console.error('保存个人信息失败:', err);
    alert('保存个人信息失败，请稍后重试');
  } finally {
    saving.value = false;
  }
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
/* 顶部导入/导出/创建按钮区域 */
.import-export-bar {
  padding: 16px 16px 10px 16px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 12px;
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
  min-width: 180px;
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
  padding: 24px 26px;
  background: linear-gradient(135deg, #fafafa 0%, #f5f7fb 100%);
}

/* 党员信息卡片样式 */
.member-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0,0,0,0.06);
  overflow: hidden;
  max-width: 980px;
  margin: 10px auto 0;
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
  padding: 26px 30px 30px;
}

.info-section {
  display: flex;
  gap: 40px;
  position: relative;
  padding-right: 180px; /* 为右侧图片区域预留空间，避免遮挡 */
}

.info-left, .info-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 15px;
  min-width: 0; /* 防止内容溢出 */
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
  right: 30px; /* 与卡片内边距对齐 */
  top: 30px;
  width: 140px; /* 稍微增加宽度 */
  height: 180px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  z-index: 1; /* 确保在内容之上但不遮挡 */
}

.photo-wrapper {
  width: 100%;
  height: 100%;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0,0,0,0.15);
  border: 1px solid #e0e0e0;
  background: #fafafa;
}

.photo-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
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

.photo-actions {
  margin-top: 6px;
}

.photo-actions .btn {
  width: 100%;
  padding: 6px 0;
  font-size: 13px;
}

@media (max-width: 1200px) {
  .container {
    width: 100%;
    margin: 0 16px;
  }

  .info-section {
    flex-direction: column;
    padding-right: 0; /* 小屏幕时移除右侧内边距 */
  }

  .photo-section {
    position: static;
    align-self: center;
    margin-top: 20px; /* 增加与上方内容的间距 */
    margin-bottom: 12px;
  }
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