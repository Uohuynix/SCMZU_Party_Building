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
            <h2>培训信息</h2>
          </div>
          <div class="actions">
            <button class="btn">导入Excel</button>
            <button class="btn">导出Excel</button>
            <button class="btn btn-primary" @click="handleAddTraining">新增</button>
          </div>
          <div class="table-wrap">
            <div v-if="loading" class="loading-state">
              加载培训数据中...
            </div>
            <div v-else-if="error" class="error-state">
              {{ error }}
            </div>
            <table v-else class="train-table">
              <thead>
                <tr>
                  <th>培训期数</th>
                  <th>姓名</th>
                  <th>民族</th>
                  <th>培训单位</th>
                  <th>培训时间</th>
                  <th>考核成绩</th>
                  <th>证书编号</th>
                  <th>操作</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="training in trainings" :key="training.id">
                  <td>{{ training.period }}</td>
                  <td><span class="highlight">{{ training.name }}</span></td>
                  <td>{{ training.ethnicity }}</td>
                  <td>{{ training.trainingUnit }}</td>
                  <td>{{ training.trainingTime }}</td>
                  <td><span class="tag" :class="getScoreTagClass(training.score)">{{ training.score }}</span></td>
                  <td>{{ training.certificateNumber }}</td>
                  <td><button class="btn btn-edit" @click="handleEditButtonClick(training)">编辑</button></td>
                </tr>
                <tr v-if="trainings.length === 0">
                  <td colspan="8" class="empty-row">暂无培训数据</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
    
    <!-- 编辑弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>培训记录操作</h3>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <p>您正在操作培训记录：</p>
          <div class="training-info">
            <p><strong>姓名：</strong>{{ currentTraining.name }}</p>
            <p><strong>培训期数：</strong>{{ currentTraining.period }}</p>
            <p><strong>当前考核成绩：</strong>{{ currentTraining.score }}</p>
            
            <!-- 更新表单 -->
            <div v-if="updateMode" class="update-form">
              <label for="newScore">新的考核成绩：</label>
              <input 
                id="newScore" 
                v-model="newScoreValue" 
                type="text" 
                placeholder="请输入新的考核成绩"
              >
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button 
            v-if="updateMode" 
            class="btn btn-cancel" 
            @click="updateMode = false"
          >
            取消更新
          </button>
          <button 
            v-else 
            class="btn btn-edit" 
            @click="prepareUpdate"
          >
            更新信息
          </button>
          <button 
            class="btn btn-danger" 
            @click="confirmDelete"
          >
            删除记录
          </button>
        </div>
      </div>
    </div>
</template>

<script setup>
import axios from 'axios';
import { ref, onMounted } from 'vue';

// API基础URL
const API_BASE_URL = 'http://127.0.0.1:4523/m1/7299957-7028709-default/api/v1';

// 响应式数据
const trainings = ref([]);
const loading = ref(false);
const error = ref('');
const showModal = ref(false);
const currentTraining = ref(null);
const updateMode = ref(false);
const newScoreValue = ref('');

// 获取培训数据
const fetchTrainings = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const response = await axios.get(`${API_BASE_URL}/trainings`);
    console.log('API返回的原始数据:', response.data);
    
    // 映射API返回的数据到表格所需格式
    trainings.value = mapTrainingData(response.data);
    console.log('映射后的培训数据:', trainings.value);
    console.log('获取培训列表成功');
  } catch (err) {
    console.error('获取培训数据失败:', err);
    error.value = '获取培训数据失败，请稍后重试';
    // 提供模拟数据以便在API调用失败时仍能显示界面
    trainings.value = getMockTrainingData();
  } finally {
    loading.value = false;
  }
};

// 映射培训数据
const mapTrainingData = (data) => {
  console.log('进入数据映射函数，输入数据类型:', typeof data);
  
  // 处理不同可能的API返回格式
  let items = [];
  
  // 如果数据是数组，直接使用
  if (Array.isArray(data)) {
    items = data;
  }
  // 如果数据是对象且有data或items属性，尝试获取
  else if (typeof data === 'object' && data !== null) {
    if (Array.isArray(data.data)) {
      items = data.data;
      console.log('从data属性获取到数组数据');
    } else if (Array.isArray(data.items)) {
      items = data.items;
      console.log('从items属性获取到数组数据');
    } else {
      // 如果对象中没有明显的数组属性，将整个对象作为单个条目
      items = [data];
      console.log('将单个对象作为数据条目');
    }
  }
  
  console.log('处理后的数据条目数量:', items.length);
  
  // 映射每个条目到表格所需格式
  return items.map((item, index) => {
    // 尝试不同可能的字段名（考虑API可能使用不同的命名约定）
    return {
      id: item.id || item.ID || index + 1,
      period: item.period || item.periodName || item.term || `2024年第${index + 1}期`,
      name: item.name || item.userName || item.personName || `党员${index + 1}`,
      ethnicity: item.ethnicity || item.nation || item.ethnic || '汉族',
      trainingUnit: item.trainingUnit || item.unit || item.organization || `党校${String.fromCharCode(65 + (index % 3))}`,
      trainingTime: item.trainingTime || item.time || item.dateRange || '2024-01-01 至 2024-01-10',
      score: item.score || item.grade || item.result || '良好',
      certificateNumber: item.certificateNumber || item.certNo || item.certificate || `NO2024${String(index + 1).padStart(3, '0')}`
    };
  });
};

// 获取成绩对应的标签样式类
const getScoreTagClass = (score) => {
  if (score === '优秀') return 'tag-red';
  if (score === '良好') return 'tag-green';
  if (score === '及格') return 'tag-yellow';
  return '';
};

// 删除培训记录
const deleteTraining = async (trainingId) => {
  try {
    const response = await axios.delete(`${API_BASE_URL}/trainings/${trainingId}`);
    console.log('删除培训记录成功:', response.data);
    return response.data;
  } catch (err) {
    console.error('删除培训记录失败:', err);
    throw err;
  }
};

// 更新培训记录
const updateTraining = async (trainingId, updatedData) => {
  try {
    const response = await axios.put(`${API_BASE_URL}/trainings/${trainingId}`, updatedData);
    console.log('更新培训记录成功:', response.data);
    return response.data;
  } catch (err) {
    console.error('更新培训记录失败:', err);
    throw err;
  }
};

// 处理编辑按钮点击事件
const handleEditButtonClick = (training) => {
  // 设置当前培训记录
  currentTraining.value = { ...training };
  updateMode.value = false;
  newScoreValue.value = training.score;
  // 显示弹窗
  showModal.value = true;
};

// 关闭弹窗
const closeModal = () => {
  showModal.value = false;
  currentTraining.value = null;
  updateMode.value = false;
  newScoreValue.value = '';
};

// 准备更新操作
const prepareUpdate = () => {
  updateMode.value = true;
};

// 确认删除操作
const confirmDelete = () => {
  if (confirm(`确定要删除培训记录「${currentTraining.value.name} - ${currentTraining.value.period}」吗？`)) {
    handleDeleteTraining(currentTraining.value);
    closeModal();
  }
};

// 处理更新培训记录
const handleUpdateTraining = async (training) => {
  // 从弹窗中获取新成绩
  const newScore = newScoreValue.value.trim();
  
  if (!newScore) {
    alert('请输入考核成绩！');
    return;
  }
  
  try {
    // 准备更新数据
    const updatedData = {
      ...training,
      score: newScore,
      updatedAt: new Date().toISOString()
    };
    
    // 调用API更新记录
    const result = await updateTraining(training.id, updatedData);
    
    // 更新本地数据
    const index = trainings.value.findIndex(item => item.id === training.id);
    if (index !== -1) {
      trainings.value[index] = mapTrainingData([result])[0];
    }
    
    alert('培训记录更新成功！');
    closeModal();
  } catch (error) {
    // API调用失败时，直接更新本地数据
    const index = trainings.value.findIndex(item => item.id === training.id);
    if (index !== -1) {
      trainings.value[index].score = newScore;
    }
    alert('培训记录已本地更新（API调用失败）！');
    closeModal();
  }
};

// 处理删除培训记录
const handleDeleteTraining = async (training) => {
  try {
    // 调用API删除记录
    await deleteTraining(training.id);
    
    // 从本地数据中移除
    const index = trainings.value.findIndex(item => item.id === training.id);
    if (index !== -1) {
      trainings.value.splice(index, 1);
    }
    
    alert('培训记录删除成功！');
  } catch (error) {
    // API调用失败时，直接从本地移除
    const index = trainings.value.findIndex(item => item.id === training.id);
    if (index !== -1) {
      trainings.value.splice(index, 1);
    }
    alert('培训记录已本地删除（API调用失败）！');
  }
};

// 添加新培训记录
const addTraining = async (trainingData) => {
  try {
    const response = await axios.post(`${API_BASE_URL}/trainings`, trainingData);
    console.log('新增培训记录成功:', response.data);
    return response.data;
  } catch (err) {
    console.error('新增培训记录失败:', err);
    throw err;
  }
};

// 处理新增按钮点击事件
const handleAddTraining = async () => {
  // 创建一个示例培训数据（实际项目中可能会从表单获取）
  const newTraining = {
    period: `2024年第${trainings.value.length + 1}期`,
    name: `新党员${trainings.value.length + 1}`,
    ethnicity: '汉族',
    trainingUnit: '党校A',
    trainingTime: new Date().toISOString().split('T')[0] + ' 至 ' + new Date(Date.now() + 7*24*60*60*1000).toISOString().split('T')[0],
    score: '良好',
    certificateNumber: `NO${new Date().getFullYear()}${String(trainings.value.length + 1).padStart(4, '0')}`
  };
  
  try {
    // 调用API新增培训记录
    const result = await addTraining(newTraining);
    
    // 将新增的记录添加到表格中
    const mappedTraining = mapTrainingData([result])[0];
    trainings.value.unshift(mappedTraining);
    
    alert('培训记录新增成功！');
  } catch (error) {
    // API调用失败时，直接将模拟数据添加到表格中
    const mappedTraining = mapTrainingData([newTraining])[0];
    trainings.value.unshift(mappedTraining);
    alert('培训记录已本地添加（API调用失败）！');
  }
};

// 模拟培训数据，当API调用失败时使用
const getMockTrainingData = () => [
  {
    id: 1,
    period: '2024年第1期',
    name: '王五',
    ethnicity: '汉族',
    trainingUnit: '党校A',
    trainingTime: '2024-05-01 至 2024-05-10',
    score: '优秀',
    certificateNumber: 'NO20240501'
  },
  {
    id: 2,
    period: '2024年第2期',
    name: '李四',
    ethnicity: '满族',
    trainingUnit: '党校B',
    trainingTime: '2024-04-15 至 2024-04-20',
    score: '良好',
    certificateNumber: 'NO20240415'
  },
  {
    id: 3,
    period: '2024年第3期',
    name: '张三',
    ethnicity: '回族',
    trainingUnit: '党校C',
    trainingTime: '2024-03-20 至 2024-03-25',
    score: '优秀',
    certificateNumber: 'NO20240320'
  },
  {
    id: 4,
    period: '2024年第4期',
    name: '赵六',
    ethnicity: '壮族',
    trainingUnit: '党校A',
    trainingTime: '2024-02-10 至 2024-02-15',
    score: '及格',
    certificateNumber: 'NO20240210'
  }
];

// 获取单个培训详情
const fetchTrainingDetail = async (trainingId = 1) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/trainings/${trainingId}`);
    console.log('获取培训详情成功:', response.data);
    
    // 如果获取到单个培训详情，将其添加到培训列表中（避免重复）
    const detailData = response.data;
    if (detailData) {
      // 检查是否已存在相同ID的培训记录
      const existingIndex = trainings.value.findIndex(item => 
        item.id === detailData.id || item.id === trainingId
      );
      
      // 如果不存在，则添加到列表中
      if (existingIndex === -1) {
        const mappedDetail = mapTrainingData([detailData])[0];
        trainings.value.unshift(mappedDetail); // 添加到列表开头
        console.log('已将培训详情添加到表格中');
      }
    }
    
    return detailData;
  } catch (err) {
    console.error('获取培训详情失败:', err);
    return null;
  }
};

// 组件挂载时获取培训数据和单个培训详情
onMounted(async () => {
  await fetchTrainings();
  console.log('培训列表加载完成，当前数据条数:', trainings.value.length);
  
  // 获取ID为1的培训详情（可以根据需要修改ID）
  await fetchTrainingDetail();
  console.log('培训详情加载完成，当前表格数据条数:', trainings.value.length);
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
.main-content {
  flex: 1;
  padding-left: 220px; /* Adjust based on sidebar width */
  width: calc(100% - 220px); /* Adjust based on sidebar width */
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
  min-width: 800px;
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
.actions {
  margin-bottom: 14px;
  display: flex;
  gap: 10px;
}
.btn {
  background: #fff;
  color: #b71c1c;
  border: 1px solid #e0e0e0;
  border-radius: 6px;
  padding: 5px 16px;
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
.btn-edit {
  background: #fff;
  color: #b71c1c;
  border: 1px solid #b71c1c;
  border-radius: 6px;
  padding: 4px 12px;
  font-size: 13px;
  margin: 0 2px;
}
.btn-edit:hover {
  background: #b71c1c;
  color: #fff;
}
.table-wrap {
  overflow-x: auto;
}
.train-table {
  width: 100%;
  border-collapse: collapse;
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  font-size: 14px;
}
.train-table th, .train-table td {
  border: 1px solid #f0f0f0;
  padding: 7px 5px;
  text-align: center;
}
.train-table th {
  background: #f5f5f5;
  color: #b71c1c;
  font-weight: bold;
}
.train-table td {
  background: #fff;
}
.highlight {
  color: #b71c1c;
  font-weight: bold;
  font-size: 15px;
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
.tag-green {
  background: #e8f5e8;
  color: #2e7d32;
  border: 1px solid #a5d6a7;
}
.tag-yellow {
  background: #fff9c4;
  color: #f57f17;
  border: 1px solid #ffd54f;
}

/* 加载状态样式 */
.loading-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  font-size: 16px;
  color: #666;
}

/* 错误状态样式 */
.error-state {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 200px;
  font-size: 16px;
  color: #e53935;
}

/* 空行样式 */
.empty-row {
  text-align: center;
  color: #999;
  padding: 30px;
  font-style: italic;
}

/* 弹窗样式 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
  border: 1px solid #e0e0e0;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px 15px;
  border-bottom: 1px solid #f0f0f0;
}

.modal-header h3 {
  margin: 0;
  color: #b71c1c;
  font-size: 1.1rem;
  font-weight: bold;
}

.modal-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
  width: 30px;
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 50%;
  transition: all 0.2s;
}

.modal-close:hover {
  background: #f5f5f5;
  color: #666;
}

.modal-body {
  padding: 20px 24px;
}

.training-info {
  background: #fafafa;
  padding: 15px;
  border-radius: 8px;
  margin: 10px 0;
  border: 1px solid #f0f0f0;
}

.training-info p {
  margin: 8px 0;
  color: #333;
}

.update-form {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px dashed #e0e0e0;
}

.update-form label {
  display: block;
  margin-bottom: 8px;
  color: #555;
  font-weight: 500;
}

.update-form input {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
}

.update-form input:focus {
  outline: none;
  border-color: #b71c1c;
  box-shadow: 0 0 0 2px rgba(183, 28, 28, 0.1);
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  padding: 15px 24px 20px;
  border-top: 1px solid #f0f0f0;
}

.btn-cancel {
  background: #f5f5f5;
  color: #666;
  border: 1px solid #ddd;
}

.btn-cancel:hover {
  background: #e0e0e0;
  color: #333;
}

.btn-danger {
  background: #e53935;
  color: white;
  border: none;
}

.btn-danger:hover {
  background: #c62828;
  color: white;
}
</style>