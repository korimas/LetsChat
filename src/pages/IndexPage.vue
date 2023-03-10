<template>
  <q-page style="min-height: 0;">

    <div class="q-pa-md column no-wrap">

      <!--聊天内容-->
      <q-scroll-area style="height: calc(100vh - 148px);" ref="scrollAreaRef" @scroll="autoScroll">
        <div class="row justify-center" >
          <div style="width: 100%; max-width: 800px;">
            <q-chat-message
              style="white-space: pre-wrap;"
              v-for="(msg, index) in Messages" :key="index"
              :name='msg.sent ? "Me": "AI"'
              :text=[msg.text]
              :avatar='msg.sent ? meImg: aiImg'
              :sent=msg.sent
              text-html
            />

            <q-chat-message
              v-if="Loading"
              name="AI"
              :avatar="aiImg"
            >
              <q-spinner-dots size="2rem"/>
            </q-chat-message>

          </div>
        </div>
      </q-scroll-area>

      <!--输入框-->
      <div class="row justify-center">
        <div class="row justify-center" style="width: 100%; max-width: 800px">
          <q-input
            square
            class="col-12 col-md-10"
            :disable="Loading"
            @keydown.enter="handleEnter"
            filled autogrow bg-color="grey"
            v-model="InputText"
            label="向OpenAI提问"/>

          <q-btn
            square
            @click="sendMessage"
            unelevated
            class="col-12 col-md-2"
            color="secondary"
          >
            <div>发送(Ctrl+Enter)</div>
          </q-btn>
        </div>
      </div>
    </div>
  </q-page>
</template>

<script lang="ts">
import {defineComponent, ref} from 'vue';
import api from "src/api/request";

type Message = {
  text: string;
  sent: boolean;
}

type GptMessage = {
  role: string;
  content: string
}

export default defineComponent({
  name: 'IndexPage',
  setup() {
    let Messages = ref<Message[]>([])
    let InputText = ref('')
    let TotalMessages = ref<GptMessage[]>([])
    let Loading = ref(false)

    const scrollAreaRef = ref()

    let scrollSize = 0
    let scrollPercent = 1

    const meImg = './imgs/me.jpg'
    const aiImg = './imgs/ai.png'

    Messages.value.push({
      sent: false,
      text: "你好，我是OpenAI小助手，基于gpt-3.5-turbo模型，采用ServerLess部署。<br>使用过程中有任何问题可联系：zpzhou.ok@gmail.com"
    })

    function sendMessage() {
      if (InputText.value == "") {
        return
      }

      Messages.value.push({
        sent: true,
        text: InputText.value
      })

      TotalMessages.value.push({
        role: "user",
        content: InputText.value
      })
      InputText.value = ""

      Loading.value = true
      api.SendMessage({
        "model": "gpt-3.5-turbo",
        "messages": TotalMessages.value
      }).then(response => {
        Loading.value = false

        // 检查出错
        console.log(response.data.result.error.type)
        if (response.data.result.error.type != "") {
          Messages.value.push({
            sent: false,
            text: "抱歉，OpenAI服务器繁忙，错误：" + response.data.result["error"]["type"]
          })
          TotalMessages.value.pop()

        } else {
          let respMessage = response.data.result.choices[0]["message"]["content"]
          TotalMessages.value?.push({
            role: "assistant",
            content: respMessage
          })

          Messages.value.push({
            sent: false,
            text: respMessage
          })
        }
      })
    }

    function scrollToBottom() {
      scrollAreaRef.value.setScrollPercentage( 'vertical', 1 )
    }

    function autoScroll() {
      const scroller = scrollAreaRef.value.getScroll()
      if (scrollSize != scroller.verticalSize) {
        // 滚动区域发生变化，判断之前是否在最底下
        if (scrollPercent === 1) {
          scrollToBottom()
        }
      }
      scrollSize = scroller.verticalSize
      scrollPercent = scroller.verticalPercentage
    }

    function handleEnter(e: any) {
      if (e.ctrlKey) {
        sendMessage()
      }
    }

    return {
      handleEnter,
      sendMessage,
      autoScroll,
      scrollAreaRef,
      InputText,
      Messages,
      Loading,
      meImg,
      aiImg
    }
  },
});
</script>