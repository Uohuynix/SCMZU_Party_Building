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
            <input
              ref="fileInput"
              type="file"
              accept=".xlsx,.xls"
              style="display:none"
              @change="handleImport"
            />
            <button v-if="canEdit" class="btn" @click="triggerImport">导入Excel</button>
            <button v-if="canEdit" class="btn" @click="exportExcel">导出Excel</button>
            <button v-if="canEdit" class="btn btn-primary" @click="handleAddTraining">新增</button>
          </div>
          <div class="table-layout">
            <!-- 左侧：按培训期数分组的筛选列 -->
            <aside class="period-panel" v-if="periodOptions.length">
              <div class="period-header">培训期数</div>
              <ul class="period-list">
                <li
                  class="period-item"
                  :class="{ active: !selectedPeriod }"
                  @click="selectedPeriod = ''"
                >
                  <span>全部期数</span>
                  <span class="period-count">{{ trainings.length }}</span>
                </li>
                <li
                  v-for="p in periodOptions"
                  :key="p.value"
                  class="period-item"
                  :class="{ active: selectedPeriod === p.value }"
                  @click="selectedPeriod = p.value"
                >
                  <span class="period-name">{{ p.label }}</span>
                  <span class="period-count">{{ p.count }}</span>
                </li>
              </ul>
            </aside>

            <!-- 右侧：培训信息表格 -->
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
                  <tr v-for="training in filteredTrainings" :key="training.id">
                    <td>{{ training.period }}</td>
                    <td><span class="highlight">{{ training.name }}</span></td>
                    <td>{{ training.ethnicity }}</td>
                    <td>{{ training.trainingUnit }}</td>
                    <td>{{ training.trainingTime }}</td>
                    <td><span class="tag" :class="getScoreTagClass(training.score)">{{ training.score }}</span></td>
                    <td>{{ training.certificateNumber }}</td>
                    <td>
                      <button v-if="canEdit" class="btn btn-edit" @click="handleEditButtonClick(training)">编辑</button>
                      <span v-else class="read-only-hint">仅查看</span>
                    </td>
                  </tr>
                  <tr v-if="filteredTrainings.length === 0">
                    <td colspan="8" class="empty-row">该期数暂无培训数据</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
    
    <!-- 新增培训信息弹窗 -->
    <div v-if="showAddModal" class="modal-overlay" @click.self="closeAddModal">
      <div class="modal-content" style="max-width: 600px;">
        <div class="modal-header">
          <h3>新增培训记录</h3>
          <button class="modal-close" @click="closeAddModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="add-form">
            <div class="form-row">
              <label for="addPeriod">培训期数 <span class="required">*</span>：</label>
              <input
                id="addPeriod"
                v-model="newTrainingForm.period"
                type="text"
                placeholder="例如：2024年第1期"
                required
              >
            </div>
            
            <div class="form-row">
              <label for="addName">姓名：</label>
              <input
                id="addName"
                v-model="newTrainingForm.name"
                type="text"
                placeholder="请输入姓名"
              >
            </div>
            
            <div class="form-row">
              <label for="addEthnicity">民族：</label>
              <input
                id="addEthnicity"
                v-model="newTrainingForm.ethnicity"
                type="text"
                placeholder="例如：汉族"
              >
            </div>
            
            <div class="form-row">
              <label for="addUnit">培训单位 <span class="required">*</span>：</label>
              <input
                id="addUnit"
                v-model="newTrainingForm.unit"
                type="text"
                placeholder="例如：党校A"
                required
              >
            </div>
            
            <div class="form-row">
              <label for="addStartDate">开始日期 <span class="required">*</span>：</label>
              <input
                id="addStartDate"
                v-model="newTrainingForm.startDate"
                type="date"
                required
              >
            </div>
            
            <div class="form-row">
              <label for="addEndDate">结束日期 <span class="required">*</span>：</label>
              <input
                id="addEndDate"
                v-model="newTrainingForm.endDate"
                type="date"
                required
              >
            </div>
            
            <div class="form-row">
              <label for="addScore">考核成绩 <span class="required">*</span>：</label>
              <select
                id="addScore"
                v-model="newTrainingForm.score"
                required
              >
                <option value="">请选择</option>
                <option value="excellent">优秀</option>
                <option value="good">良好</option>
                <option value="pass">及格</option>
                <option value="fail">不及格</option>
              </select>
            </div>
            
            <div class="form-row">
              <label for="addCertificateNo">证书编号：</label>
              <input
                id="addCertificateNo"
                v-model="newTrainingForm.certificateNo"
                type="text"
                placeholder="例如：NO20240001"
              >
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-cancel" @click="closeAddModal">取消</button>
          <button class="btn btn-primary" @click="handleSaveNewTraining">保存</button>
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
            <p><strong>当前证书编号：</strong>{{ currentTraining.certificateNumber || '（暂无）' }}</p>
            
            <!-- 更新表单 -->
            <div v-if="updateMode" class="update-form">
              <label for="newPeriod">新的培训期数：</label>
              <input
                id="newPeriod"
                v-model="newPeriodValue"
                type="text"
                placeholder="请输入新的培训期数"
              >
              <label for="newScore" style="margin-top:10px;">新的考核成绩：</label>
              <input 
                id="newScore" 
                v-model="newScoreValue" 
                type="text" 
                placeholder="请输入新的考核成绩"
              >
              <label for="newCert" style="margin-top:10px;">新的证书编号：</label>
              <input
                id="newCert"
                v-model="newCertificateValue"
                type="text"
                placeholder="如不修改可留空"
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
          <button
            v-if="updateMode"
            class="btn btn-primary"
            @click="handleUpdateTraining(currentTraining)"
          >
            保存修改
          </button>
        </div>
      </div>
    </div>
</template>

<script setup>
import axios from 'axios';
import * as XLSX from 'xlsx';
import { ref, onMounted, computed } from 'vue';

// API基础URL（通过 Vite 代理到本地后端）
const API_BASE_URL = '/api/v1';
const token = localStorage.getItem('token') || '';

// 获取用户角色并判断是否有编辑权限
const userRole = localStorage.getItem('role') || 'student';
const canEdit = ref(userRole === 'admin' || userRole === 'teacher');

// 响应式数据
const trainings = ref([]);
const loading = ref(false);
const error = ref('');
const showModal = ref(false);
const currentTraining = ref(null);
const updateMode = ref(false);
const newScoreValue = ref('');
const newPeriodValue = ref('');
const newCertificateValue = ref('');

// 按培训期数分组筛选
const selectedPeriod = ref('');
const periodOptions = computed(() => {
  const map = new Map();
  trainings.value.forEach(t => {
    const key = t.period || '未填写期数';
    if (!map.has(key)) {
      map.set(key, { value: key, label: key, count: 0 });
    }
    map.get(key).count += 1;
  });
  return Array.from(map.values());
});

const filteredTrainings = computed(() => {
  if (!selectedPeriod.value) return trainings.value;
  return trainings.value.filter(t => (t.period || '未填写期数') === selectedPeriod.value);
});

// 新增培训信息相关
const showAddModal = ref(false);
const newTrainingForm = ref({
  period: '',
  name: '',
  ethnicity: '',
  unit: '',
  startDate: '',
  endDate: '',
  score: '',
  certificateNo: ''
});

// Excel 导入用的隐藏文件输入
const fileInput = ref(null);

// 获取培训数据
const fetchTrainings = async () => {
  loading.value = true;
  error.value = '';
  
  try {
    const response = await axios.get(`${API_BASE_URL}/trainings`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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

// 后端分数枚举 <-> 中文文案 映射
const scoreEnumToLabel = (value) => {
  const map = {
    excellent: '优秀',
    good: '良好',
    pass: '及格',
    fail: '不及格'
  };
  return map[value] || value || '';
};

const scoreLabelToEnum = (label) => {
  if (!label) return '';
  
  // 去除空格
  label = label.trim();
  
  // 定义映射关系
  const map = {
    '优秀': 'excellent',
    '良好': 'good',
    '及格': 'pass',
    '不及格': 'fail'
  };
  
  // 先尝试完全匹配
  if (map[label]) {
    return map[label];
  }
  
  // 如果完全匹配失败，尝试从输入中提取有效的中文成绩
  // 例如："优秀1" -> "优秀"
  for (const [chinese, enumValue] of Object.entries(map)) {
    if (label.includes(chinese)) {
      return enumValue;
    }
  }
  
  // 如果都不匹配，尝试直接使用原值（可能是枚举值）
  const validEnums = ['excellent', 'good', 'pass', 'fail'];
  if (validEnums.includes(label.toLowerCase())) {
    return label.toLowerCase();
  }
  
  // 如果都不匹配，返回空字符串（后端会验证）
  return '';
};

// 映射培训数据（专门针对当前后端返回结构）
const mapTrainingData = (raw) => {
  const items = Array.isArray(raw) ? raw : (raw?.data || []);

  return items.map((item, index) => {
    const start = item.start_date ? String(item.start_date).split('T')[0] : '';
    const end = item.end_date ? String(item.end_date).split('T')[0] : '';
    const timeRange = start && end ? `${start} 至 ${end}` : '';

    const student = item.student || {};

    return {
      id: item.id || index + 1,
      period: item.period || '',
      name: student.name || '',
      ethnicity: student.ethnicity || '',
      trainingUnit: item.unit || '',
      trainingTime: timeRange,
      score: scoreEnumToLabel(item.score),
      certificateNumber: item.certificate_no || '',
      // 供更新 / 导出 / 导入时使用的原始字段
      _studentId: item.student_id || null,
      _startDate: start,
      _endDate: end,
      _scoreEnum: item.score || ''
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
    const response = await axios.delete(`${API_BASE_URL}/trainings/${trainingId}`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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
    const response = await axios.put(`${API_BASE_URL}/trainings/${trainingId}`, updatedData, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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
  newPeriodValue.value = training.period;
  newCertificateValue.value = training.certificateNumber || '';
  // 显示弹窗
  showModal.value = true;
};

// 关闭弹窗
const closeModal = () => {
  showModal.value = false;
  currentTraining.value = null;
  updateMode.value = false;
  newScoreValue.value = '';
  newPeriodValue.value = '';
  newCertificateValue.value = '';
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
  // 从弹窗中获取新培训期数和新成绩
  const newScore = newScoreValue.value.trim();
  const newPeriod = newPeriodValue.value.trim();
  const newCert = newCertificateValue.value.trim();
  
  if (!newPeriod) {
    alert('请输入培训期数！');
    return;
  }
  if (!newScore) {
    alert('请输入考核成绩！');
    return;
  }

  const enumScore = scoreLabelToEnum(newScore);
  
  // 验证转换后的成绩是否有效
  const validEnums = ['excellent', 'good', 'pass', 'fail'];
  if (!enumScore || !validEnums.includes(enumScore)) {
    alert('请输入有效的考核成绩：优秀、良好、及格或不及格');
    return;
  }
  
  const time = (training.trainingTime || '').split('至').map(s => s.trim());
  let start = time[0] || training._startDate || null;
  let end = time[1] || training._endDate || null;
  
  // 确保日期格式正确（只取日期部分，去掉时间部分）
  if (start && start.includes(' ')) {
    start = start.split(' ')[0];
  }
  if (end && end.includes(' ')) {
    end = end.split(' ')[0];
  }
  
  try {
    // 准备更新数据
    const updatedData = {
      period: newPeriod,
      student_id: training._studentId || undefined,
      unit: training.trainingUnit || '',
      start_date: start || '',
      end_date: end || '',
      score: enumScore,
      certificate_no: newCert || training.certificateNumber || ''
    };
    
    // 移除undefined字段
    Object.keys(updatedData).forEach(key => {
      if (updatedData[key] === undefined) {
        delete updatedData[key];
      }
    });
    
    // 调用API更新记录
    const result = await updateTraining(training.id, updatedData);

    // 直接整体刷新列表，防止本地状态与后端不一致
    await fetchTrainings();

    alert('培训记录更新成功！');
    closeModal();
  } catch (error) {
    console.error('更新培训记录失败:', error);
    let errorMessage = '培训记录更新失败，请稍后重试！';
    if (error.response && error.response.data && error.response.data.error) {
      errorMessage = `更新失败: ${error.response.data.error}`;
    } else if (error.message) {
      errorMessage = `更新失败: ${error.message}`;
    }
    alert(errorMessage);
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
    const response = await axios.post(`${API_BASE_URL}/trainings`, trainingData, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
    console.log('新增培训记录成功:', response.data);
    return response.data;
  } catch (err) {
    console.error('新增培训记录失败:', err);
    throw err;
  }
};

// 处理新增按钮点击事件
const handleAddTraining = () => {
  // 重置表单
  newTrainingForm.value = {
    period: '',
    name: '',
    ethnicity: '',
    unit: '',
    startDate: '',
    endDate: '',
    score: '',
    certificateNo: ''
  };
  // 显示新增弹窗
  showAddModal.value = true;
};

// 关闭新增弹窗
const closeAddModal = () => {
  showAddModal.value = false;
  // 重置表单
  newTrainingForm.value = {
    period: '',
    name: '',
    ethnicity: '',
    unit: '',
    startDate: '',
    endDate: '',
    score: '',
    certificateNo: ''
  };
};

// 保存新增的培训记录
const handleSaveNewTraining = async () => {
  // 验证必填字段
  if (!newTrainingForm.value.period) {
    alert('请输入培训期数！');
    return;
  }
  if (!newTrainingForm.value.unit) {
    alert('请输入培训单位！');
    return;
  }
  if (!newTrainingForm.value.startDate) {
    alert('请选择开始日期！');
    return;
  }
  if (!newTrainingForm.value.endDate) {
    alert('请选择结束日期！');
    return;
  }
  if (!newTrainingForm.value.score) {
    alert('请选择考核成绩！');
    return;
  }
  
  // 验证日期逻辑
  if (new Date(newTrainingForm.value.startDate) > new Date(newTrainingForm.value.endDate)) {
    alert('结束日期不能早于开始日期！');
    return;
  }

  try {
    // 查找或创建学生（如果有姓名）
    let studentId = null;
    if (newTrainingForm.value.name) {
      try {
        // 获取所有学生列表进行查找
        const studentsResponse = await axios.get(`${API_BASE_URL}/students`, {
          params: { page: 1, pageSize: 1000 },
          headers: {
            Authorization: token ? `Bearer ${token}` : ''
          }
        });
        
        let student = null;
        if (studentsResponse.data && studentsResponse.data.data && Array.isArray(studentsResponse.data.data)) {
          // 精确匹配姓名
          student = studentsResponse.data.data.find(s => s.name && s.name.trim() === newTrainingForm.value.name.trim());
        }
        
        if (!student) {
          // 如果学生不存在，创建新学生
          const newStudent = {
            student_no: `AUTO_${Date.now()}`,
            name: newTrainingForm.value.name,
            ethnicity: newTrainingForm.value.ethnicity || '汉族',
            type: 'masses',
            gender: 'male'
          };
          
          const createResponse = await axios.post(`${API_BASE_URL}/students`, newStudent, {
            headers: {
              Authorization: token ? `Bearer ${token}` : ''
            }
          });
          
          if (createResponse.data && createResponse.data.id) {
            studentId = createResponse.data.id;
          }
        } else {
          studentId = student.id;
          // 如果民族信息不同，更新学生信息
          if (newTrainingForm.value.ethnicity && student.ethnicity !== newTrainingForm.value.ethnicity) {
            try {
              await axios.put(`${API_BASE_URL}/students/${student.id}`, {
                ethnicity: newTrainingForm.value.ethnicity
              }, {
                headers: {
                  Authorization: token ? `Bearer ${token}` : ''
                }
              });
            } catch (updateErr) {
              console.warn('更新学生民族信息失败:', updateErr);
            }
          }
        }
      } catch (studentErr) {
        console.warn('处理学生信息失败:', studentErr);
        // 继续创建培训记录，但不关联学生
      }
    }

    // 准备培训数据
    const newTraining = {
      period: newTrainingForm.value.period,
      unit: newTrainingForm.value.unit,
      start_date: newTrainingForm.value.startDate,
      end_date: newTrainingForm.value.endDate,
      score: newTrainingForm.value.score,
      certificate_no: newTrainingForm.value.certificateNo || ''
    };
    
    if (studentId) {
      newTraining.student_id = studentId;
    }

    // 调用API新增培训记录
    await addTraining(newTraining);
    await fetchTrainings();

    alert('培训记录新增成功！');
    closeAddModal();
  } catch (error) {
    console.error('新增培训记录失败:', error);
    let errorMessage = '新增培训记录失败，请稍后重试！';
    if (error.response && error.response.data && error.response.data.error) {
      errorMessage = `新增失败: ${error.response.data.error}`;
    } else if (error.message) {
      errorMessage = `新增失败: ${error.message}`;
    }
    alert(errorMessage);
  }
};

// 触发 Excel 导入
const triggerImport = () => {
  if (fileInput.value) fileInput.value.value = '';
  fileInput.value?.click();
};

// Excel日期转换函数
const excelDateToJSDate = (excelDate) => {
  if (!excelDate && excelDate !== 0) return null;
  
  // 如果是数字（Excel日期序列号）
  if (typeof excelDate === 'number') {
    try {
      // Excel日期从1900-01-01开始，但Excel错误地认为1900是闰年
      const excelEpoch = new Date(1899, 11, 30);
      const jsDate = new Date(excelEpoch.getTime() + excelDate * 24 * 60 * 60 * 1000);
      const result = jsDate.toISOString().split('T')[0];
      console.log(`Excel日期序列号 ${excelDate} 转换为: ${result}`);
      return result;
    } catch (e) {
      console.warn('Excel日期序列号转换失败:', excelDate, e);
      return null;
    }
  }
  
  // 如果是字符串
  const str = String(excelDate).trim();
  if (!str || str === 'null' || str === 'undefined') return null;
  
  // 如果已经是YYYY-MM-DD格式
  if (/^\d{4}-\d{2}-\d{2}$/.test(str)) {
    return str;
  }
  
  // 尝试解析为日期对象
  try {
    const date = new Date(str);
    if (!isNaN(date.getTime())) {
      const result = date.toISOString().split('T')[0];
      console.log(`日期字符串 "${str}" 转换为: ${result}`);
      return result;
    }
  } catch (e) {
    console.warn('日期解析失败:', str, e);
  }
  
  // 尝试匹配常见日期格式
  const datePatterns = [
    /(\d{4})[-\/](\d{1,2})[-\/](\d{1,2})/,  // YYYY-MM-DD 或 YYYY/MM/DD
    /(\d{2})[-\/](\d{1,2})[-\/](\d{4})/,   // MM-DD-YYYY 或 MM/DD/YYYY
  ];
  
  for (const pattern of datePatterns) {
    const match = str.match(pattern);
    if (match) {
      let year, month, day;
      if (match[1].length === 4) {
        // YYYY-MM-DD格式
        year = match[1];
        month = match[2].padStart(2, '0');
        day = match[3].padStart(2, '0');
      } else {
        // MM-DD-YYYY格式
        month = match[1].padStart(2, '0');
        day = match[2].padStart(2, '0');
        year = match[3];
      }
      const result = `${year}-${month}-${day}`;
      console.log(`日期格式匹配 "${str}" 转换为: ${result}`);
      return result;
    }
  }
  
  console.warn('无法解析日期:', str);
  return null;
};

// Excel 导入处理
const handleImport = async (e) => {
  const file = e.target.files?.[0];
  if (!file) {
    alert('请选择要导入的Excel文件');
    return;
  }

  const reader = new FileReader();
  reader.onload = async (evt) => {
    try {
      console.log('开始读取Excel文件...');
      const data = new Uint8Array(evt.target.result);
      const workbook = XLSX.read(data, { 
        type: 'array', 
        cellDates: false,  // 不自动转换日期，手动处理
        cellNF: false,
        raw: false
      });
      
      if (!workbook.SheetNames || workbook.SheetNames.length === 0) {
        alert('Excel文件没有工作表');
        return;
      }
      
      const sheet = workbook.Sheets[workbook.SheetNames[0]];
      console.log('读取工作表:', workbook.SheetNames[0]);
      
      // 直接使用sheet_to_json，自动识别表头
      const rows = XLSX.utils.sheet_to_json(sheet, { 
        defval: '',
        raw: false,
        dateNF: 'yyyy-mm-dd'
      });
      
      console.log('读取到', rows.length, '行数据');
      console.log('第一行数据示例:', rows[0]);
      
      if (rows.length === 0) {
        alert('Excel文件中没有数据');
        return;
      }

      let successCount = 0;
      let failCount = 0;
      const errors = [];

      // 按行导入
      for (let i = 0; i < rows.length; i++) {
        const row = rows[i];
        
        // 读取所有字段（支持多种可能的列名，不区分大小写）
        const getField = (fieldNames) => {
          for (const name of fieldNames) {
            // 尝试精确匹配
            if (row[name] !== undefined && row[name] !== null && row[name] !== '') {
              return String(row[name]).trim();
            }
            // 尝试不区分大小写匹配
            const lowerName = name.toLowerCase();
            for (const key in row) {
              if (key.toLowerCase() === lowerName) {
                return String(row[key]).trim();
              }
            }
          }
          return '';
        };
        
        const period = getField(['培训期数', '期数', 'period', 'Period']);
        const name = getField(['姓名', 'name', 'Name']);
        const ethnicity = getField(['民族', 'ethnicity', 'Ethnicity']);
        const unit = getField(['培训单位', '单位', 'unit', 'Unit']);
        const timeStr = getField(['培训时间', '时间', 'time', 'Time']);
        const scoreLabel = getField(['考核成绩', '成绩', 'score', 'Score']);
        const certNo = getField(['证书编号', '证书', 'certificate_no', 'CertificateNo']);
        
        // 尝试读取分开的开始和结束日期
        let startDateStr = getField(['开始日期', 'start_date', 'startDate', 'StartDate']);
        let endDateStr = getField(['结束日期', 'end_date', 'endDate', 'EndDate']);
        
        console.log(`处理第${i + 2}行:`, { period, name, unit, scoreLabel });

      // 验证必填字段
      if (!period || !unit || !scoreLabel) {
        failCount++;
        errors.push(`第${i + 2}行：缺少必填字段（培训期数、培训单位、考核成绩）`);
        continue;
      }

      // 解析日期
      let start = '';
      let end = '';
      
      // 优先使用分开的开始和结束日期列
      if (startDateStr && endDateStr) {
        start = excelDateToJSDate(startDateStr) || startDateStr;
        end = excelDateToJSDate(endDateStr) || endDateStr;
      } else if (timeStr) {
        // 如果时间列包含"至"，分割处理
        if (timeStr.includes('至')) {
          const timeParts = timeStr.split('至').map(s => s.trim());
          start = excelDateToJSDate(timeParts[0]) || timeParts[0] || '';
          end = excelDateToJSDate(timeParts[1]) || timeParts[1] || '';
        } else {
          // 单个日期，尝试解析
          const parsed = excelDateToJSDate(timeStr);
          if (parsed) {
            start = parsed;
            end = parsed;
          }
        }
      }
      
      // 清理日期格式（确保是YYYY-MM-DD格式）
      if (start) {
        // 移除可能的时间部分
        start = start.split(' ')[0].split('T')[0];
        const startDate = new Date(start);
        if (!isNaN(startDate.getTime())) {
          start = startDate.toISOString().split('T')[0];
        } else {
          // 尝试其他格式
          const parts = start.match(/(\d{4})[-\/](\d{1,2})[-\/](\d{1,2})/);
          if (parts) {
            start = `${parts[1]}-${parts[2].padStart(2, '0')}-${parts[3].padStart(2, '0')}`;
          } else {
            start = ''; // 如果无法解析，清空
          }
        }
      }
      
      if (end) {
        // 移除可能的时间部分
        end = end.split(' ')[0].split('T')[0];
        const endDate = new Date(end);
        if (!isNaN(endDate.getTime())) {
          end = endDate.toISOString().split('T')[0];
        } else {
          // 尝试其他格式
          const parts = end.match(/(\d{4})[-\/](\d{1,2})[-\/](\d{1,2})/);
          if (parts) {
            end = `${parts[1]}-${parts[2].padStart(2, '0')}-${parts[3].padStart(2, '0')}`;
          } else {
            end = ''; // 如果无法解析，清空
          }
        }
      }

      // 转换成绩
      const scoreEnum = scoreLabelToEnum(scoreLabel);
      if (!scoreEnum) {
        failCount++;
        errors.push(`第${i + 2}行：考核成绩格式不正确（应为：优秀、良好、及格、不及格）`);
        continue;
      }

      // 查找或创建学生（如果有姓名）
      let studentId = null;
      if (name) {
        try {
          // 获取所有学生列表进行查找（对于导入场景，教师和管理员可以查看所有学生）
          const studentsResponse = await axios.get(`${API_BASE_URL}/students`, {
            params: { page: 1, pageSize: 1000 }, // 获取足够多的学生
            headers: {
              Authorization: token ? `Bearer ${token}` : ''
            }
          });
          
          let student = null;
          if (studentsResponse.data && studentsResponse.data.data && Array.isArray(studentsResponse.data.data)) {
            // 精确匹配姓名
            student = studentsResponse.data.data.find(s => s.name && s.name.trim() === name);
          }
          
          if (!student) {
            // 如果学生不存在，创建新学生
            const newStudent = {
              student_no: `AUTO_${Date.now()}_${i}`, // 自动生成学号
              name: name,
              ethnicity: ethnicity || '汉族',
              type: 'masses', // 默认为群众
              gender: 'male' // 默认性别
            };
            
            try {
              const createResponse = await axios.post(`${API_BASE_URL}/students`, newStudent, {
                headers: {
                  Authorization: token ? `Bearer ${token}` : ''
                }
              });
              
              if (createResponse.data && createResponse.data.id) {
                studentId = createResponse.data.id;
                console.log(`创建新学生: ${name}, ID: ${studentId}`);
              }
            } catch (createErr) {
              console.warn(`第${i + 2}行：创建学生失败:`, createErr);
              // 继续导入，但不关联学生
            }
          } else {
            studentId = student.id;
            // 如果民族信息不同，更新学生信息
            if (ethnicity && student.ethnicity !== ethnicity) {
              try {
                await axios.put(`${API_BASE_URL}/students/${student.id}`, {
                  ethnicity: ethnicity
                }, {
                  headers: {
                    Authorization: token ? `Bearer ${token}` : ''
                  }
                });
              } catch (updateErr) {
                console.warn('更新学生民族信息失败:', updateErr);
              }
            }
          }
        } catch (studentErr) {
          console.warn(`第${i + 2}行：处理学生信息失败:`, studentErr);
          // 继续导入，但不关联学生
        }
      }

      // 准备培训数据
      const payload = {
        period: period,
        unit: unit,
        score: scoreEnum,
        certificate_no: certNo
      };
      
      if (start) payload.start_date = start;
      if (end) payload.end_date = end;
      if (studentId) payload.student_id = studentId;

        try {
          console.log(`导入第${i + 2}行数据:`, payload);
          await addTraining(payload);
          successCount++;
          console.log(`第${i + 2}行导入成功`);
        } catch (err) {
          failCount++;
          const errorMsg = err.response?.data?.error || err.message || '未知错误';
          errors.push(`第${i + 2}行：${errorMsg}`);
          console.error(`导入第${i + 2}行失败：`, {
            row: row,
            payload: payload,
            error: err.response?.data || err.message
          });
        }
      }

      console.log('导入完成，刷新数据...');
      await fetchTrainings();
      
      // 显示导入结果
      let message = `Excel 导入完成！\n成功：${successCount} 条`;
      if (failCount > 0) {
        message += `\n失败：${failCount} 条`;
        if (errors.length > 0) {
          message += `\n\n错误详情（前10条）：\n${errors.slice(0, 10).join('\n')}`;
          if (errors.length > 10) {
            message += `\n...还有${errors.length - 10}条错误`;
          }
        }
      }
      alert(message);
    } catch (importError) {
      console.error('Excel导入过程出错:', importError);
      alert(`Excel导入失败：${importError.message || '未知错误'}\n\n请检查：\n1. Excel文件格式是否正确\n2. 是否包含表头行\n3. 必填字段是否完整\n4. 日期格式是否正确`);
    }
  };
  reader.readAsArrayBuffer(file);
};

// Excel 导出
const exportExcel = () => {
  if (!trainings.value.length) {
    alert('当前没有可导出的培训数据');
    return;
  }

  // 导出数据，确保格式与导入兼容
  const rows = trainings.value.map(t => {
    // 解析培训时间，提取开始和结束日期
    let startDate = '';
    let endDate = '';
    if (t.trainingTime) {
      const timeParts = t.trainingTime.split('至').map(s => s.trim());
      startDate = timeParts[0] || '';
      endDate = timeParts[1] || '';
    }
    
    // 如果原始数据中有日期，优先使用
    if (t._startDate) startDate = t._startDate;
    if (t._endDate) endDate = t._endDate;
    
    return {
      '培训期数': t.period || '',
      '姓名': t.name || '',
      '民族': t.ethnicity || '',
      '培训单位': t.trainingUnit || '',
      '培训时间': t.trainingTime || (startDate && endDate ? `${startDate} 至 ${endDate}` : ''),
      '开始日期': startDate || '',
      '结束日期': endDate || '',
      '考核成绩': t.score || '',
      '证书编号': t.certificateNumber || ''
    };
  });

  // 创建工作表
  const ws = XLSX.utils.json_to_sheet(rows);
  
  // 设置列宽
  const colWidths = [
    { wch: 15 }, // 培训期数
    { wch: 10 }, // 姓名
    { wch: 8 },  // 民族
    { wch: 15 }, // 培训单位
    { wch: 25 }, // 培训时间
    { wch: 12 }, // 开始日期
    { wch: 12 }, // 结束日期
    { wch: 10 }, // 考核成绩
    { wch: 15 }  // 证书编号
  ];
  ws['!cols'] = colWidths;
  
  const wb = XLSX.utils.book_new();
  XLSX.utils.book_append_sheet(wb, ws, '培训信息');
  
  // 生成文件名（包含时间戳）
  const fileName = `培训信息_${new Date().toISOString().split('T')[0]}.xlsx`;
  XLSX.writeFile(wb, fileName);
  
  console.log('Excel导出成功，共导出', rows.length, '条记录');
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
    const response = await axios.get(`${API_BASE_URL}/trainings/${trainingId}`, {
      headers: {
        Authorization: token ? `Bearer ${token}` : ''
      }
    });
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

/* 表格整体布局：左侧期数列表 + 右侧表格 */
.table-layout {
  display: flex;
  margin-top: 8px;
}

.period-panel {
  width: 220px;
  border-right: 1px solid #f0f0f0;
  padding-right: 10px;
  margin-right: 10px;
}

.period-header {
  font-size: 14px;
  font-weight: 600;
  color: #b71c1c;
  margin-bottom: 6px;
}

.period-list {
  list-style: none;
  margin: 0;
  padding: 0;
  max-height: 420px;
  overflow-y: auto;
}

.period-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 6px 10px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: #555;
  margin-bottom: 4px;
  transition: background 0.2s, color 0.2s;
}

.period-item:hover {
  background: #fff5f5;
  color: #b71c1c;
}

.period-item.active {
  background: #b71c1c;
  color: #fff;
}

.period-name {
  flex: 1;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.period-count {
  font-size: 12px;
  margin-left: 8px;
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

.read-only-hint {
  color: #999;
  font-size: 13px;
  font-style: italic;
}

/* 新增表单样式 */
.add-form {
  padding: 10px 0;
}

.form-row {
  margin-bottom: 20px;
}

.form-row label {
  display: block;
  margin-bottom: 8px;
  color: #333;
  font-weight: 500;
  font-size: 14px;
}

.form-row .required {
  color: #e53935;
  margin-left: 2px;
}

.form-row input[type="text"],
.form-row input[type="date"],
.form-row select {
  width: 100%;
  padding: 10px 12px;
  border: 1px solid #ddd;
  border-radius: 6px;
  font-size: 14px;
  transition: border-color 0.2s;
  box-sizing: border-box;
}

.form-row input[type="text"]:focus,
.form-row input[type="date"]:focus,
.form-row select:focus {
  outline: none;
  border-color: #b71c1c;
  box-shadow: 0 0 0 2px rgba(183, 28, 28, 0.1);
}

.form-row select {
  background-color: #fff;
  cursor: pointer;
}
</style>