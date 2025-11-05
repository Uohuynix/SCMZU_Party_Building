<template>
  <div class="dashboard-bg">
    <header class="dashboard-header">
      <img src="/src/pictures/dang.png" alt="党徽" class="logo" />
      <h1>智慧党建信息管理平台</h1>
    </header>
    <div class="dashboard-content">
      <aside class="sidebar">
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
      <main class="main-panel">
        <!-- 党校校长视图 - 查看全部党支部 -->
        <div v-if="userRole === '党校校长'" class="dean-view">
          <div class="stats-row">
            <div class="stat-card">
              <div class="stat-title">正式党员</div>
              <div class="stat-num stat-red">{{ allStats.formal }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-title">预备党员</div>
              <div class="stat-num">{{ allStats.probationary }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-title">发展对象</div>
              <div class="stat-num">{{ allStats.developing }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-title">入党积极分子</div>
              <div class="stat-num">{{ allStats.activist }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-title">入党申请人</div>
              <div class="stat-num">{{ allStats.applicant }}</div>
            </div>
          </div>
                   <div class="charts-row">
            <div class="chart-card">
              <div class="chart-title">性别占比</div>
              <div class="gender-pie-chart">
                <div class="pie-chart">
                  <div class="pie-center">
                    <div class="pie-total">总计</div>
                    <div class="pie-number">{{ genderData.total }}</div>
                  </div>
                </div>
                <div class="pie-legend">
                  <div class="legend-item">
                    <div class="legend-color male"></div>
                    <span>男: {{ genderData.male }} ({{ genderData.malePercent }}%)</span>
                  </div>
                  <div class="legend-item">
                    <div class="legend-color female"></div>
                    <span>女: {{ genderData.female }} ({{ genderData.femalePercent }}%)</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="chart-card">
              <div class="chart-title">年龄分布</div>
              <div class="age-horizontal-chart">
                <div class="age-horizontal-bars">
                  <div v-for="item in ageRangeData" :key="item.range" class="age-horizontal-bar">
                    <div class="age-horizontal-label">{{ item.range }}</div>
                    <div class="age-horizontal-bar-container">
                      <div class="age-horizontal-fill" :style="{ width: item.percent + '%' }"></div>
                    </div>
                    <div class="age-horizontal-count">{{ item.count }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="charts-row">
            <div class="chart-card">
              <div class="chart-title">民族结构</div>
              <div class="ethnic-vertical-chart">
                <div class="ethnic-vertical-bars">
                  <div v-for="item in ethnicBarData" :key="item.name" class="ethnic-vertical-bar">
                    <div class="ethnic-vertical-fill" :style="{ height: item.percent + '%' }"></div>
                    <div class="ethnic-vertical-count">{{ item.count }}</div>
                    <div class="ethnic-vertical-name">{{ item.name }}</div>
                  </div>
                </div>
              </div>
            </div>
            <div class="chart-card">
              <div class="chart-title">学历结构</div>
              <div class="education-radar-chart">
                <div class="radar-chart">
                  <svg width="150" height="150" viewBox="0 0 150 150">
                    <polygon
                      :points="radarPoints"
                      fill="rgba(231, 76, 60, 0.3)"
                      stroke="#e74c3c"
                      stroke-width="2"
                    />
                    <circle cx="75" cy="75" r="3" fill="#e74c3c" />
                  </svg>
                </div>
              </div>
            </div>
          </div>

         <div class="notice-row">
           <div class="notice-card">
             <div class="notice-title">近期发展动态/通知公告</div>
             <ul class="notice-list">
               <li v-for="notice in notices" :key="notice.date">{{ notice.date }} {{ notice.content }}</li>
             </ul>
           </div>
         </div>
        </div>

        <!-- 书记/副书记视图 - 只能查看本年级 -->
        <div v-else class="secretary-view">
          <div class="stats-row">
            <div class="stat-card">
              <div class="stat-title">正式党员</div>
              <div class="stat-num stat-red">{{ branchStats.formal }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-title">预备党员</div>
              <div class="stat-num">{{ branchStats.probationary }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-title">发展党员</div>
              <div class="stat-num">{{ branchStats.developing }}</div>
            </div>
          </div>
          <div class="charts-row">
            <div class="chart-card">
              <div class="chart-title">性别占比</div>
              <div class="gender-pie-chart">
                <div class="pie-chart">
                  <div class="pie-center">
                    <div class="pie-total">总计</div>
                    <div class="pie-number">{{ branchGenderData.total }}</div>
                  </div>
                </div>
                <div class="pie-legend">
                  <div class="legend-item">
                    <div class="legend-color male"></div>
                    <span>男: {{ branchGenderData.male }} ({{ branchGenderData.malePercent }}%)</span>
                  </div>
                  <div class="legend-item">
                    <div class="legend-color female"></div>
                    <span>女: {{ branchGenderData.female }} ({{ branchGenderData.femalePercent }}%)</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="chart-card">
              <div class="chart-title">年龄分布</div>
              <div class="age-horizontal-chart">
                <div class="age-horizontal-bars">
                  <div v-for="item in branchAgeRangeData" :key="item.range" class="age-horizontal-bar">
                    <div class="age-horizontal-label">{{ item.range }}</div>
                    <div class="age-horizontal-bar-container">
                      <div class="age-horizontal-fill" :style="{ width: item.percent + '%' }"></div>
                    </div>
                    <div class="age-horizontal-count">{{ item.count }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="charts-row">
            <div class="chart-card">
              <div class="chart-title">民族结构</div>
              <div class="ethnic-vertical-chart">
                <div class="ethnic-vertical-bars">
                  <div v-for="item in branchEthnicBarData" :key="item.name" class="ethnic-vertical-bar">
                    <div class="ethnic-vertical-fill" :style="{ height: item.percent + '%' }"></div>
                    <div class="ethnic-vertical-count">{{ item.count }}</div>
                    <div class="ethnic-vertical-name">{{ item.name }}</div>
                  </div>
                </div>
              </div>
            </div>
            <div class="chart-card">
              <div class="chart-title">学历结构</div>
              <div class="education-radar-chart">
                <div class="radar-chart">
                  <svg width="150" height="150" viewBox="0 0 150 150">
                    <polygon
                      :points="radarPoints"
                      fill="rgba(231, 76, 60, 0.3)"
                      stroke="#e74c3c"
                      stroke-width="2"
                    />
                    <circle cx="75" cy="75" r="3" fill="#e74c3c" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
          <div class="notice-row">
            <div class="notice-card">
              <div class="notice-title">近期发展动态/通知公告</div>
              <ul class="notice-list">
                <li v-for="notice in notices.slice(0, 3)" :key="notice.date">{{ notice.date }} {{ notice.content }}</li>
              </ul>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import axios from 'axios';

// 获取用户角色和相关信息
const userRole = localStorage.getItem('role') || '党校校长';
const userBranch = localStorage.getItem('branch') || 'all';
const userGrade = localStorage.getItem('grade') || 'all';

// 党校校长视图的统计数据（全校）
const allStats = ref({
  formal: 31,
  probationary: 29,
  developing: 44,
  activist: 55,
  applicant: 60
});

// 书记/副书记视图的统计数据（本年级）
const branchStats = ref({
  formal: 45,
  probationary: 35,
  developing: 20,
  partyBranch: 4
});

// 性别统计数据
const genderData = ref({
  total: 120,
  male: 65,
  female: 55,
  malePercent: 54,
  femalePercent: 46
});

// 年龄范围数据
const ageRangeData = ref([
  { range: '20-24岁', count: 35, percent: 29 },
  { range: '25-34岁', count: 45, percent: 38 },
  { range: '35-44岁', count: 25, percent: 21 },
  { range: '45-59岁', count: 12, percent: 10 },
  { range: '≥60岁', count: 3, percent: 2 }
]);

// 民族条形图数据
const ethnicBarData = ref([
  { name: '汉族', count: 95, percent: 79 },
  { name: '满族', count: 8, percent: 7 },
  { name: '回族', count: 6, percent: 5 },
  { name: '壮族', count: 4, percent: 3 },
  { name: '苗族', count: 3, percent: 2 },
  { name: '其他', count: 4, percent: 4 }
]);

// 通知公告数据
const notices = ref([
  { date: '2024-05-01', content: '全校党员大会顺利召开' },
  { date: '2024-04-20', content: '各党支部换届选举完成' },
  { date: '2024-04-10', content: '全校党员培训圆满结束' },
  { date: '2024-03-15', content: '研究生党支部成立仪式' }
]);

// 计算年龄函数
const calculateAge = (birthDateStr) => {
  const birthDate = new Date(birthDateStr);
  const today = new Date();
  let age = today.getFullYear() - birthDate.getFullYear();
  const monthDiff = today.getMonth() - birthDate.getMonth();
  if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birthDate.getDate())) {
    age--;
  }
  return age;
};

// 获取党员类型对应的统计类别
const getPartyTypeCategory = (type) => {
  switch (type) {
    case 'formal': return 'formal';
    case 'probationary': return 'probationary';
    case 'developing': return 'developing';
    case 'activist': return 'activist';
    case 'applicant': return 'applicant';
    default: return 'masses';
  }
};

// 获取年龄段
const getAgeRange = (age) => {
  if (age >= 20 && age <= 24) return '20-24岁';
  if (age >= 25 && age <= 34) return '25-34岁';
  if (age >= 35 && age <= 44) return '35-44岁';
  if (age >= 45 && age <= 59) return '45-59岁';
  if (age >= 60) return '≥60岁';
  return '其他';
};

// 加载数据的函数
const loadData = async () => {
  try {
    // 模拟获取多个党员数据（假设我们需要获取多个党员信息）
    const studentIds = [1, 2, 3, 4, 5]; // 假设需要获取这些党员的信息
    const studentsData = [];
    
    // 逐个获取党员信息
    for (const id of studentIds) {
      const response = await axios.get(`http://127.0.0.1:4523/m1/7299957-7028709-default/api/v1/students/${id}`);
      if (response.data) {
        studentsData.push(response.data);
      }
    }
    
    // 如果获取到数据，更新统计信息
    if (studentsData.length > 0) {
      // 初始化统计计数器
      const statsCounter = {
        formal: 0,
        probationary: 0,
        developing: 0,
        activist: 0,
        applicant: 0
      };
      
      const genderCounter = { male: 0, female: 0 };
      const ethnicCounter = {};
      const ageRangeCounter = {};
      const educationCounter = {};
      
      // 统计数据
      studentsData.forEach(student => {
        // 统计党员类型
        const category = getPartyTypeCategory(student.type);
        if (statsCounter[category] !== undefined) {
          statsCounter[category]++;
        }
        
        // 统计性别
        if (student.gender === 'male') genderCounter.male++;
        if (student.gender === 'female') genderCounter.female++;
        
        // 统计民族
        ethnicCounter[student.ethnicity] = (ethnicCounter[student.ethnicity] || 0) + 1;
        
        // 统计年龄
        const age = calculateAge(student.birth_date);
        const ageRange = getAgeRange(age);
        ageRangeCounter[ageRange] = (ageRangeCounter[ageRange] || 0) + 1;
        
        // 统计学历
        educationCounter[student.education] = (educationCounter[student.education] || 0) + 1;
      });
      
      // 更新全校统计数据
      allStats.value = {
        formal: statsCounter.formal,
        probationary: statsCounter.probationary,
        developing: statsCounter.developing,
        activist: statsCounter.activist,
        applicant: statsCounter.applicant
      };
      
      // 更新党支部统计数据（这里简化处理，实际应该根据branch和grade筛选）
      branchStats.value = {
        formal: Math.floor(statsCounter.formal * 0.7),
        probationary: Math.floor(statsCounter.probationary * 0.7),
        developing: Math.floor(statsCounter.developing * 0.7),
        partyBranch: 4
      };
      
      // 更新性别统计
      const totalGender = genderCounter.male + genderCounter.female;
      genderData.value = {
        total: totalGender,
        male: genderCounter.male,
        female: genderCounter.female,
        malePercent: totalGender > 0 ? Math.round((genderCounter.male / totalGender) * 100) : 0,
        femalePercent: totalGender > 0 ? Math.round((genderCounter.female / totalGender) * 100) : 0
      };
      
      // 更新年龄分布
      const totalAge = studentsData.length;
      ageRangeData.value = Object.entries(ageRangeCounter).map(([range, count]) => ({
        range,
        count,
        percent: Math.round((count / totalAge) * 100)
      })).sort((a, b) => {
        // 按年龄范围排序
        const rangeOrder = ['20-24岁', '25-34岁', '35-44岁', '45-59岁', '≥60岁', '其他'];
        return rangeOrder.indexOf(a.range) - rangeOrder.indexOf(b.range);
      });
      
      // 更新民族结构
      const totalEthnic = studentsData.length;
      ethnicBarData.value = Object.entries(ethnicCounter)
        .map(([name, count]) => ({
          name,
          count,
          percent: Math.round((count / totalEthnic) * 100)
        }))
        .sort((a, b) => b.count - a.count)
        .slice(0, 6); // 只显示前6个民族
      
      // 更新党支部性别数据
      branchGenderData.value = {
        total: Math.floor(totalGender * 0.7),
        male: Math.floor(genderCounter.male * 0.7),
        female: Math.floor(genderCounter.female * 0.7),
        malePercent: genderData.value.malePercent,
        femalePercent: genderData.value.femalePercent
      };
      
      // 更新党支部年龄数据
      branchAgeRangeData.value = ageRangeData.value.map(item => ({
        ...item,
        count: Math.floor(item.count * 0.7),
        percent: item.percent
      }));
      
      // 更新党支部民族数据
      branchEthnicBarData.value = ethnicBarData.value.slice(0, 4).map(item => ({
        ...item,
        count: Math.floor(item.count * 0.7)
      }));
      
      console.log('数据更新成功，基于API获取的党员信息');
    }
    
  } catch (error) {
    console.error('获取数据失败:', error);
    // 保持原有模拟数据作为备用
  }
};

// 组件挂载时加载数据
onMounted(() => {
  loadData();
});



// 雷达图数据点 - 学历结构
const radarPoints = ref('80,15 125,80 80,125 15,80 45,15 115,15');

// 书记/副书记视图数据
const branchGenderData = ref({
  total: 80,
  male: 42,
  female: 38,
  malePercent: 53,
  femalePercent: 47
});

const branchAgeRangeData = ref([
  { range: '20-24岁', count: 25, percent: 31 },
  { range: '25-34岁', count: 35, percent: 44 },
  { range: '35-44岁', count: 15, percent: 19 },
  { range: '45-59岁', count: 4, percent: 5 },
  { range: '≥60岁', count: 1, percent: 1 }
]);

const branchEthnicBarData = ref([
  { name: '汉族', count: 68, percent: 85 },
  { name: '满族', count: 6, percent: 8 },
  { name: '回族', count: 4, percent: 5 },
  { name: '其他', count: 2, percent: 2 }
]);


</script>

<style scoped>
.dashboard-bg {
  min-height: 100vh;
  background: linear-gradient(135deg, #fbe9e7 60%, #b71c1c 100%);
}
.dashboard-header {
  display: flex;
  align-items: center;
  padding: 28px 40px 16px 40px;
  background: rgba(255,255,255,0.12);
  border-bottom: 1px solid #e57373;
  box-shadow: 0 2px 12px rgba(183,28,28,0.08);
}
.logo {
  width: 48px;
  height: 48px;
  margin-right: 18px;
  border-radius: 50%;
  border: 2px solid #b71c1c;
  background: #fff;
}
h1 {
  color: #b71c1c;
  font-size: 2.1rem;
  font-weight: bold;
  letter-spacing: 2px;
  text-shadow: 1px 2px 8px #fff5f5;
}
.dashboard-content {
  display: flex;
  height: calc(100vh - 80px);
}
.sidebar {
  width: 180px;
  background: rgba(255,255,255,0.18);
  border-right: 1px solid #e57373;
  padding: 32px 0 0 0;
  display: flex;
  flex-direction: column;
  gap: 18px;
}
.sidebar a {
  color: #b71c1c;
  font-weight: 500;
  font-size: 16px;
  padding: 10px 24px;
  border-radius: 8px 0 0 8px;
  text-decoration: none;
  transition: background 0.2s, color 0.2s;
  display: flex;
  align-items: center;
  gap: 12px;
}
.sidebar-icon {
  width: 20px;
  height: 20px;
  border-radius: 50%;
  border: 1px solid #b71c1c;
  background: #fff;
}
.sidebar a.router-link-exact-active, .sidebar a:hover {
  background: #fff5f5;
  color: #b71c1c;
}
.main-panel {
  flex: 1;
  padding: 25px 30px 0 30px;
  overflow-y: auto;
}
.stats-row {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 15px;
  margin-bottom: 25px;
}
.stat-card {
  background: linear-gradient(135deg, #fff 0%, #f8f9fa 100%);
  border-radius: 12px;
  box-shadow: 0 3px 15px rgba(243,156,18,0.12);
  padding: 15px 20px;
  display: flex;
  flex-direction: column;
  align-items: center;
  border: 2px solid #f39c12;
  position: relative;
  overflow: hidden;
}

.stat-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 3px;
  background: linear-gradient(90deg, #f39c12 0%, #e67e22 100%);
}
.stat-title {
  color: #e74c3c;
  font-size: 12px;
  margin-bottom: 6px;
  letter-spacing: 0.5px;
}
.stat-num {
  color: #f39c12;
  font-size: 1.8rem;
  font-weight: bold;
  letter-spacing: 0.5px;
  text-shadow: 0 1px 4px rgba(243,156,18,0.2);
}
.stat-red {
  color: #e74c3c;
}
.charts-row {
  display: flex;
  gap: 20px;
  margin-bottom: 25px;
}

.center-chart {
  justify-content: center;
}

.center-chart .chart-card {
  max-width: 400px;
}
.chart-card {
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  padding: 20px;
  flex: 1 1 0;
  min-width: 280px;
  border: 1px solid #f0f0f0;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.chart-title {
  color: #e74c3c;
  font-size: 16px;
  font-weight: 600;
  margin-bottom: 15px;
  letter-spacing: 0.5px;
}
.chart-placeholder {
  width: 100%;
  height: 160px;
  background: #fce4ec;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #e74c3c;
  font-size: 14px;
  font-weight: 500;
  margin-top: 10px;
}
.notice-row {
  margin-bottom: 32px;
}
.notice-card {
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 2px 12px rgba(183,28,28,0.08);
  padding: 26px 38px;
  border: 1.5px solid #fbe9e7;
}
.notice-title {
  color: #b71c1c;
  font-size: 16px;
  font-weight: bold;
  margin-bottom: 10px;
  letter-spacing: 1px;
}
.notice-list {
  list-style: none;
  padding: 0;
  margin: 0;
  color: #616161;
  font-size: 15px;
}
.notice-list li {
  padding: 3px 0;
}

/* 性别饼状图样式 */
.gender-pie-chart {
  padding: 20px;
  text-align: center;
}

.pie-chart {
  position: relative;
  width: 100px;
  height: 100px;
  margin: 0 auto 12px;
}

.pie-center {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  text-align: center;
  background: white;
  border-radius: 50%;
  width: 50px;
  height: 50px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.pie-total {
  font-size: 12px;
  color: #666;
}

.pie-number {
  font-size: 14px;
  font-weight: bold;
  color: #e74c3c;
}

.pie-legend {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 15px;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #666;
}

.legend-color {
  width: 16px;
  height: 16px;
  border-radius: 2px;
}

.legend-color.male {
  background: #f39c12;
}

.legend-color.female {
  background: #e74c3c;
}

/* 年龄横向柱状图样式 */
.age-horizontal-chart {
  padding: 20px;
}

.age-horizontal-bars {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.age-horizontal-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.age-horizontal-label {
  width: 80px;
  font-size: 12px;
  color: #666;
}

.age-horizontal-bar-container {
  flex: 1;
  height: 16px;
  background: #f5f5f5;
  border-radius: 8px;
  overflow: hidden;
}

.age-horizontal-fill {
  height: 100%;
  background: linear-gradient(90deg, #e74c3c 0%, #f39c12 100%);
  transition: width 0.3s ease;
}

.age-horizontal-count {
  width: 40px;
  font-size: 12px;
  color: #666;
  text-align: right;
}

/* 民族垂直柱状图样式 */
.ethnic-vertical-chart {
  padding: 15px;
  height: 160px;
}

.ethnic-vertical-bars {
  display: flex;
  justify-content: space-around;
  align-items: flex-end;
  height: 100%;
  gap: 10px;
}

.ethnic-vertical-bar {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  flex: 1;
}

.ethnic-vertical-fill {
  width: 100%;
  background: linear-gradient(180deg, #e74c3c 0%, #f39c12 100%);
  border-radius: 4px 4px 0 0;
  transition: height 0.3s ease;
  min-height: 20px;
}

.ethnic-vertical-count {
  font-size: 12px;
  color: #666;
  font-weight: bold;
}

.ethnic-vertical-name {
  font-size: 10px;
  color: #666;
  text-align: center;
  writing-mode: vertical-rl;
  text-orientation: mixed;
}



/* 雷达图样式 */
.education-radar-chart {
  padding: 20px;
  text-align: center;
}

.radar-chart {
  margin: 0 auto;
}



/* 饼图样式 */
.ethnic-chart, .education-chart {
  padding: 20px;
  text-align: center;
}

.pie-chart {
  position: relative;
  width: 120px;
  height: 120px;
  margin: 0 auto 20px;
  border-radius: 50%;
  background: #f5f5f5;
}



.ethnic-legend, .education-legend {
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 12px;
  color: #666;
}

.ethnic-legend .legend-item, .education-legend .legend-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.ethnic-legend .legend-color, .education-legend .legend-color {
  width: 12px;
  height: 12px;
  border-radius: 2px;
}
</style>