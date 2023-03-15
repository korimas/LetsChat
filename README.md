# 介绍

Let's chat是一个后端接入OpenAI API的网页版聊天机器人，使用vercel提供的serverless服务部署。

前端框架使用Quasar。

[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fkorimas%2FLetsChat&env=OPENAI_API_KEY)

Demo地址：[https://chat.zpzhou.com](https://chat.zpzhou.com)


# 环境变量
部署后需要配置环境变量：
* OPENAI_API_KEY：在OpenAI账户生成的api key
* PASSWORD：保护密码，不设置时，任何人可访问
