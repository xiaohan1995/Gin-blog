/**
 * 封装 AJAX 请求方法
 */

const Ajax = {
  /**
   * 发送 GET 请求
   * @param {string} url - 请求地址
   * @param {object} params - 请求参数
   * @returns {Promise}
   */
  get(url, params = {}) {
    return this.request('GET', url, params);
  },

  /**
   * 发送 POST 请求
   * @param {string} url - 请求地址
   * @param {object} data - 请求体数据
   * @returns {Promise}
   */
  post(url, data = {}) {
    return this.request('POST', url, {}, data);
  },

  /**
   * 发送 PUT 请求
   * @param {string} url - 请求地址
   * @param {object} data - 请求体数据
   * @returns {Promise}
   */
  put(url, data = {}) {
    return this.request('PUT', url, {}, data);
  },

  /**
   * 发送 DELETE 请求
   * @param {string} url - 请求地址
   * @param {object} params - 请求参数
   * @returns {Promise}
   */
  delete(url, params = {}) {
    return this.request('DELETE', url, params);
  },

  /**
   * 通用请求方法
   * @param {string} method - 请求方法
   * @param {string} url - 请求地址
   * @param {object} params - 请求参数
   * @param {object} data - 请求体数据
   * @returns {Promise}
   */
  request(method, url, params = {}, data = {}) {
    return new Promise((resolve, reject) => {
      // 构建完整的 URL 和查询参数
      let fullUrl = url;
      if (Object.keys(params).length > 0) {
        const queryString = new URLSearchParams(params).toString();
        fullUrl += (fullUrl.includes('?') ? '&' : '?') + queryString;
      }

      // 创建 XMLHttpRequest 对象
      const xhr = new XMLHttpRequest();

      // 初始化请求
      xhr.open(method, fullUrl, true);

      // 设置默认请求头
      xhr.setRequestHeader('Content-Type', 'application/json;charset=UTF-8');
      xhr.setRequestHeader('Accept', 'application/json');

      // 获取并设置 token
      const token = localStorage.getItem('token');
      if (token) {
        xhr.setRequestHeader('Authorization', 'Bearer ' + token);
      }

      // 处理响应
      xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
          if (xhr.status >= 200 && xhr.status < 300) {
            try {
              const response = JSON.parse(xhr.responseText);
              resolve(response);
            } catch (e) {
              resolve(xhr.responseText);
            }
          } else {
            try {
              const errorResponse = JSON.parse(xhr.responseText);
              reject(new Error(errorResponse.message || '请求失败'));
            } catch (e) {
              reject(new Error('请求失败'));
            }
          }
        }
      };

      // 处理网络错误
      xhr.onerror = function () {
        reject(new Error('网络错误'));
      };

      // 发送请求
      if (method === 'POST' || method === 'PUT') {
        xhr.send(JSON.stringify(data));
      } else {
        xhr.send();
      }
    });
  }
};