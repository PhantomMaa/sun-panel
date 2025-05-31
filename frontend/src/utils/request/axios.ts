import axios, { type AxiosResponse } from 'axios'
import { useAuthStore } from '../../store'

const service = axios.create({
  baseURL: import.meta.env.VITE_GLOB_API_URL,
})

service.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()
    const token = authStore.token
    
    // 从 URL 路径中获取 public code
    let publiccode = ''
    const pathSegments = window.location.pathname.split('/')
    if (pathSegments.length > 1 && pathSegments[1] !== '') {
      publiccode = pathSegments[1]
    }
    
    // 添加 publiccode 到请求头（如果存在）
    if (publiccode)
      config.headers.publiccode = publiccode
    else
      config.headers.Authorization = `Bearer ${token}`
    
    return config
  },
  (error) => {
    return Promise.reject(error.response)
  },
)

service.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    if (response.status === 200)
      return response

    throw new Error(response.status.toString())
  },
  (error) => {
    return Promise.reject(error)
  },
)

export default service
