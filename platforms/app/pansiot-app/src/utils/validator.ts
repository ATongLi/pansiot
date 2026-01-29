/**
 * 表单验证工具
 */

/**
 * 验证手机号
 */
export const validatePhone = (phone: string): boolean => {
  const reg = /^1[3-9]\d{9}$/;
  return reg.test(phone);
};

/**
 * 验证邮箱
 */
export const validateEmail = (email: string): boolean => {
  const reg = /^[a-zA-Z0-9._-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,6}$/;
  return reg.test(email);
};

/**
 * 验证密码 (6-20位,包含字母和数字)
 */
export const validatePassword = (password: string): boolean => {
  const reg = /^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,20}$/;
  return reg.test(password);
};

/**
 * 验证用户名 (4-20位,字母开头,只能包含字母、数字、下划线)
 */
export const validateUsername = (username: string): boolean => {
  const reg = /^[a-zA-Z][a-zA-Z0-9_]{3,19}$/;
  return reg.test(username);
};

/**
 * 验证身份证号
 */
export const validateIdCard = (idCard: string): boolean => {
  const reg = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/;
  return reg.test(idCard);
};

/**
 * 验证URL
 */
export const validateUrl = (url: string): boolean => {
  const reg = /^https?:\/\/(([a-zA-Z0-9_-])+(\.)?)*(:\d+)?(\/((\.)?(\?)?=?&?[a-zA-Z0-9_-](\?)?)*)*$/i;
  return reg.test(url);
};

/**
 * 验证IP地址
 */
export const validateIp = (ip: string): boolean => {
  const reg = /^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$/;
  return reg.test(ip);
};

/**
 * 验证是否为空
 */
export const validateEmpty = (value: any): boolean => {
  if (value === null || value === undefined) return true;
  if (typeof value === 'string' && value.trim() === '') return true;
  if (Array.isArray(value) && value.length === 0) return true;
  if (typeof value === 'object' && Object.keys(value).length === 0) return true;
  return false;
};

/**
 * 验证码验证 (6位数字)
 */
export const validateCode = (code: string): boolean => {
  const reg = /^\d{6}$/;
  return reg.test(code);
};
