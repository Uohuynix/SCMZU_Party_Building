<template>
  <div class="login-bg">
    <div class="login-box">
      <h2>智慧党建信息管理平台</h2>
      <input placeholder="账号" v-model="username" />
      <input placeholder="密码" type="password" v-model="password" />
      <button @click="login">登录</button>
      <p class="error-message" v-if="errorMessage">{{ errorMessage }}</p>
      <p class="register-link">没有账号？<a href="/register">立即注册</a></p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
const username = ref('');
const password = ref('');
const errorMessage = ref('');
const router = useRouter();

function login() {
  // 清空错误消息
  errorMessage.value = '';
  
  // 验证表单
  if (!username.value || !password.value) {
    errorMessage.value = '请输入账号和密码';
    return;
  }
  
  // 准备登录数据
  const loginData = {
    username: username.value,
    password: password.value
  };
  
  // 调用登录API
  fetch('http://127.0.0.1:4523/m1/7299957-7028709-default/api/v1/auth/login', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(loginData)
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('登录失败，请检查账号密码');
    }
    return response.json();
  })
  .then(data => {
    // 登录成功，存储用户信息
    localStorage.setItem('role', data.role || 'student');
    localStorage.setItem('branch', data.branch || 'all');
    localStorage.setItem('grade', data.grade || 'all');
    localStorage.setItem('token', data.token || '');
    console.log("登陆成功");
    console.log(data);
    // 跳转到首页
    router.push('/dashboard');
  })
  .catch(error => {
    // 如果API调用失败，尝试使用本地验证作为后备
    console.log('API调用失败，尝试本地验证', error);
    localLoginFallback();
  });
}

// 本地验证后备逻辑
function localLoginFallback() {
  // 尝试从注册用户中查找
  const registeredUsers = JSON.parse(localStorage.getItem('registeredUsers') || '[]');
  const user = registeredUsers.find(u => u.username === username.value && u.password === password.value);
  
  if (user) {
    // 找到注册用户
    localStorage.setItem('role', user.role);
    localStorage.setItem('branch', user.branch);
    localStorage.setItem('grade', user.grade);
    router.push('/dashboard');
  } else {
    // 尝试使用原来的自动识别逻辑（保留兼容性）
    let role, branch, grade;

    if (username.value === 'admin' && password.value === 'admin') {
      // 党校校长 - 可以查看全部
      role = '党校校长';
      branch = 'all';
      grade = 'all';
    } else if (username.value.includes('2021')) {
      // 2021级书记/副书记
      role = '书记';
      branch = '本科生党支部';
      grade = '2021级';
    } else if (username.value.includes('2022')) {
      // 2022级书记/副书记
      role = '书记';
      branch = '本科生党支部';
      grade = '2022级';
    } else if (username.value.includes('2023')) {
      // 2023级书记/副书记
      role = '书记';
      branch = '本科生党支部';
      grade = '2023级';
    } else if (username.value.includes('2024')) {
      // 2024级书记/副书记
      role = '书记';
      branch = '本科生党支部';
      grade = '2024级';
    } else if (username.value.includes('研究生')) {
      // 研究生党支部书记/副书记
      role = '书记';
      branch = '研究生党支部';
      grade = '研究生';
    } else {
      errorMessage.value = '账号或密码错误，请重试';
      return;
    }

    localStorage.setItem('role', role);
    localStorage.setItem('branch', branch);
    localStorage.setItem('grade', grade);
    router.push('/dashboard');
  }
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
</style>