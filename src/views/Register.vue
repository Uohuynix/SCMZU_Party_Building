<template>
  <div class="login-bg">
    <div class="login-box">
      <h2>智慧党建信息管理平台 - 注册</h2>
      <input placeholder="账号" v-model="username" />
      <input placeholder="密码" type="password" v-model="password" />
      <input placeholder="确认密码" type="password" v-model="confirmPassword" />
      <input placeholder="姓名" v-model="name" />
      <select v-model="role">
        <option value="student">学生</option>
        <option value="teacher">教师</option>
        <option value="admin">管理员</option>
      </select>
      <!-- 学生需录入支部+小组；教师只需录入支部；管理员不属于任何具体党支部 -->
      <input
        v-if="role === 'teacher' || role === 'student'"
        placeholder="所属党支部（如：计算机学院学生第一党支部）"
        v-model="branch"
      />
      <input
        v-if="role === 'student'"
        placeholder="所属党小组（如：第一党小组）"
        v-model="groupName"
      />
      <button @click="register">注册</button>
      <p class="register-link">已有账号？<a href="/login">返回登录</a></p>
      <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>
      <p class="success-message" v-if="successMessage">{{ successMessage }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const username = ref('');
const password = ref('');
const confirmPassword = ref('');
const name = ref('');
const role = ref('student'); // 默认值为student
const branch = ref('');
const groupName = ref('');
const errorMessage = ref('');
const successMessage = ref('');
const router = useRouter();

function register() {
  // 清空消息
  errorMessage.value = '';
  successMessage.value = '';
  
  // 验证表单
  if (!username.value || !password.value || !confirmPassword.value || !name.value) {
    errorMessage.value = '请填写所有必填字段';
    return;
  }
  
  if (password.value !== confirmPassword.value) {
    errorMessage.value = '两次输入的密码不一致';
    return;
  }
  
  // 角色相关校验：
  // 学生：必须填写党支部和党小组
  // 教师：只需填写党支部
  if (role.value === 'student') {
    if (!branch.value || !groupName.value) {
      errorMessage.value = '学生必须填写所属党支部和党小组';
      return;
    }
  } else if (role.value === 'teacher') {
    if (!branch.value) {
      errorMessage.value = '教师必须填写所属党支部';
      return;
    }
  }

  // 准备API请求数据
  const registerData = {
    username: username.value,
    password: password.value,
    name: name.value,
    role: role.value,
    branch: role.value === 'admin' ? '' : branch.value,
    group_name: role.value === 'student' ? groupName.value : ''
  };
  
  // 调用后端注册API（通过 Vite 代理到 http://localhost:8080）
  fetch('/api/v1/auth/register', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(registerData)
  })
    .then(async response => {
      const data = await response.json().catch(() => ({}));
      if (!response.ok) {
        // 优先展示后端返回的 error 字段
        const msg = data.error || '注册失败，请稍后重试';
        throw new Error(msg);
      }
      return data;
    })
    .then(data => {
      // 注册成功
      successMessage.value = '注册成功！即将跳转到登录页面...';
      console.log(data);
      setTimeout(() => {
        router.push('/login');
      }, 1500);
    })
    .catch(error => {
      errorMessage.value = error.message || '注册失败，请稍后重试';
    });
}
</script>

<style scoped>
.login-bg {
  background: #b71c1c;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
}
.login-box {
  background: #fff;
  padding: 40px 32px;
  border-radius: 8px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.1);
  display: flex;
  flex-direction: column;
  gap: 16px;
  min-width: 320px;
  align-items: center;
}
.login-box h2 {
  color: #b71c1c;
}
.login-box input, .login-box select {
  width: 100%;
  padding: 8px;
  border: 1px solid #ccc;
  border-radius: 4px;
}
.login-box button {
  width: 100%;
  padding: 10px;
  background: #b71c1c;
  color: #fff;
  border: none;
  border-radius: 4px;
  font-size: 16px;
  cursor: pointer;
}
.register-link {
  margin-top: 8px;
  font-size: 14px;
}
.register-link a {
  color: #b71c1c;
  text-decoration: none;
}
.register-link a:hover {
  text-decoration: underline;
}
.error-message {
  color: #f44336;
  font-size: 14px;
  margin-top: -8px;
}
.success-message {
  color: #4caf50;
  font-size: 14px;
  margin-top: -8px;
}
</style>