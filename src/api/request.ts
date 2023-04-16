// import {Dialog} from 'quasar'
import type {AxiosInstance, AxiosRequestConfig, AxiosResponse} from 'axios';
import axios from 'axios';
import {Dialog} from 'quasar';

// 导出Request类，可以用来自定义传递配置来创建实例
export class Request {
  // axios 实例
  instance: AxiosInstance;
  // 基础配置，url和超时时间
  baseConfig: AxiosRequestConfig = {baseURL: '/api', timeout: 60000};

  // 类构造函数
  constructor(config: AxiosRequestConfig) {
    // 使用axios.create创建axios实例
    this.instance = axios.create(Object.assign(this.baseConfig, config));

    /* // 拦截请求
    this.instance.interceptors.request.use(
      (config: any) => {
        // 一般会请求拦截里面加token，用于后端的验证
        const token = localStorage.getItem("token") as string
        if(token) {
          config.headers!.Authorization = token;
        }

        return config;
      },
      (err: any) => {
        // 请求错误，这里可以用全局提示框进行提示
        return Promise.reject(err);
      }
    );*/

    // 拦截应答
    this.instance.interceptors.response.use(
      (res: AxiosResponse) => {
        // 直接返回res，当然你也可以只返回res.data
        // 系统如果有自定义code也可以在这里处理
        if (!res.data.success) {
          Dialog.create({
              title: '错误',
              message: res.data.message
            });
        }
        return res;
      },
      (err: any) => {
        // 这里用来处理http常见错误，进行全局提示
        Dialog.create({
          title: '错误',
          message: '请求异常，http status：' + err.response.status
        });
        // 这里是AxiosError类型，所以一般我们只reject我们需要的响应即可
        return Promise.reject(err.response);
      }
    );
  }

  // 定义请求方法
  public request(config: AxiosRequestConfig): Promise<AxiosResponse> {
    return this.instance.request(config);
  }

  public CheckNeedAuth() {
    return this.instance.get(
      '/auth',
      this.baseConfig);
  }

  public PasswordAuth(password: string) {
      return this.instance.post(
          '/auth',
          {password: password},
          this.baseConfig);
  }

}

// 默认导出Request实例
export default new Request({})
