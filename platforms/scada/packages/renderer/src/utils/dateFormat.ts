/**
 * 日期格式化工具
 * 提供日期格式化和相对时间显示功能
 */

/**
 * 时间单位定义（毫秒）
 */
const TIME_UNITS = {
  minute: 60 * 1000,
  hour: 60 * 60 * 1000,
  day: 24 * 60 * 60 * 1000,
  week: 7 * 24 * 60 * 60 * 1000,
  month: 30 * 24 * 60 * 60 * 1000,
  year: 365 * 24 * 60 * 60 * 1000
}

/**
 * 格式化相对时间
 * @param date 日期对象或ISO字符串
 * @returns 相对时间字符串（如 "2小时前"）
 */
export function formatRelativeTime(date: Date | string): string {
  const now = Date.now()
  const past = typeof date === 'string' ? new Date(date).getTime() : date.getTime()
  const diff = now - past

  // 未来时间
  if (diff < 0) {
    return '刚刚'
  }

  // 小于1分钟
  if (diff < TIME_UNITS.minute) {
    return '刚刚'
  }

  // 小于1小时
  if (diff < TIME_UNITS.hour) {
    const minutes = Math.floor(diff / TIME_UNITS.minute)
    return `${minutes}分钟前`
  }

  // 小于1天
  if (diff < TIME_UNITS.day) {
    const hours = Math.floor(diff / TIME_UNITS.hour)
    return `${hours}小时前`
  }

  // 小于1周
  if (diff < TIME_UNITS.week) {
    const days = Math.floor(diff / TIME_UNITS.day)
    return `${days}天前`
  }

  // 小于1个月
  if (diff < TIME_UNITS.month) {
    const weeks = Math.floor(diff / TIME_UNITS.week)
    return `${weeks}周前`
  }

  // 小于1年
  if (diff < TIME_UNITS.year) {
    const months = Math.floor(diff / TIME_UNITS.month)
    return `${months}月前`
  }

  // 大于1年
  const years = Math.floor(diff / TIME_UNITS.year)
  return `${years}年前`
}

/**
 * 格式化日期为 ISO 8601 字符串
 * @param date 日期对象
 * @returns ISO 8601 字符串
 */
export function formatISODate(date: Date): string {
  return date.toISOString()
}

/**
 * 格式化日期为本地字符串
 * @param date 日期对象或ISO字符串
 * @param format 格式类型 ('full' | 'long' | 'medium' | 'short')
 * @returns 格式化的日期字符串
 */
export function formatLocaleDate(
  date: Date | string,
  format: 'full' | 'long' | 'medium' | 'short' = 'medium'
): string {
  const dateObj = typeof date === 'string' ? new Date(date) : date

  const formatOptions: Record<string, Intl.DateTimeFormatOptions> = {
    full: {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    },
    long: {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    },
    medium: {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    },
    short: {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour12: false
    }
  }

  return dateObj.toLocaleString('zh-CN', formatOptions[format])
}

/**
 * 格式化日期为自定义格式
 * @param date 日期对象或ISO字符串
 * @param format 格式字符串 (如 'YYYY-MM-DD HH:mm:ss')
 * @returns 格式化的日期字符串
 */
export function formatCustomDate(
  date: Date | string,
  format: string
): string {
  const dateObj = typeof date === 'string' ? new Date(date) : date

  const year = dateObj.getFullYear()
  const month = String(dateObj.getMonth() + 1).padStart(2, '0')
  const day = String(dateObj.getDate()).padStart(2, '0')
  const hours = String(dateObj.getHours()).padStart(2, '0')
  const minutes = String(dateObj.getMinutes()).padStart(2, '0')
  const seconds = String(dateObj.getSeconds()).padStart(2, '0')

  return format
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 解析ISO字符串为Date对象
 * @param isoString ISO 8601字符串
 * @returns Date对象
 */
export function parseISODate(isoString: string): Date {
  return new Date(isoString)
}

/**
 * 获取今天的开始时间（00:00:00）
 * @returns 今天的开始时间
 */
export function getStartOfDay(): Date {
  const date = new Date()
  date.setHours(0, 0, 0, 0)
  return date
}

/**
 * 获取今天的结束时间（23:59:59）
 * @returns 今天的结束时间
 */
export function getEndOfDay(): Date {
  const date = new Date()
  date.setHours(23, 59, 59, 999)
  return date
}

/**
 * 判断是否为今天
 * @param date 日期对象或ISO字符串
 * @returns 是否为今天
 */
export function isToday(date: Date | string): boolean {
  const dateObj = typeof date === 'string' ? new Date(date) : date
  const today = new Date()

  return (
    dateObj.getFullYear() === today.getFullYear() &&
    dateObj.getMonth() === today.getMonth() &&
    dateObj.getDate() === today.getDate()
  )
}

/**
 * 判断是否为本周
 * @param date 日期对象或ISO字符串
 * @returns 是否为本周
 */
export function isThisWeek(date: Date | string): boolean {
  const dateObj = typeof date === 'string' ? new Date(date) : date
  const today = new Date()
  const startOfWeek = new Date(today)
  startOfWeek.setDate(today.getDate() - today.getDay())
  startOfWeek.setHours(0, 0, 0, 0)

  const endOfWeek = new Date(startOfWeek)
  endOfWeek.setDate(startOfWeek.getDate() + 6)
  endOfWeek.setHours(23, 59, 59, 999)

  return dateObj >= startOfWeek && dateObj <= endOfWeek
}
