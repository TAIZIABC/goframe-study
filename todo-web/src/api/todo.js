/**
 * Todo API 请求层
 * 对接后端: GET/POST/PUT/DELETE /api/v1/todos
 * 响应格式: { code: 0, message: "", data: { ... } }
 */

const BASE_URL = '/api/v1';

async function request(url, options = {}) {
  const res = await fetch(`${BASE_URL}${url}`, {
    headers: { 'Content-Type': 'application/json' },
    ...options,
  });
  const json = await res.json();
  if (json.code !== 0) {
    throw new Error(json.message || '请求失败');
  }
  return json.data;
}

/** 获取 Todo 列表 */
export async function fetchTodos(page = 1, size = 10) {
  return request(`/todos?page=${page}&size=${size}`);
}

/** 创建 Todo */
export async function createTodo(title) {
  return request('/todos', {
    method: 'POST',
    body: JSON.stringify({ title, completed: false }),
  });
}

/** 更新 Todo */
export async function updateTodo(id, data) {
  return request(`/todos/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

/** 删除 Todo */
export async function deleteTodo(id) {
  return request(`/todos/${id}`, { method: 'DELETE' });
}
