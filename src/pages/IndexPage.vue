<template>
    <q-page style="min-height: 0;">

        <div class="q-pa-md column no-wrap">

            <!--聊天内容-->
            <q-scroll-area style="height: calc(100vh - 148px);" ref="scrollAreaRef" @scroll="autoScroll">
                <div class="row justify-center">
                    <div style="width: 100%; max-width: 800px;">
                        <q-chat-message
                            style="white-space: pre-wrap;"
                            v-for="(msg, index) in DisplayMessages" :key="index"
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

                        <q-chat-message
                            style="white-space: pre-wrap;"
                            v-if="Waiting"
                            name="AI"
                            :avatar="aiImg"
                            :text=[waitText]
                        />

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
                        @click="StreamChat"
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
        let DisplayMessages = ref<Message[]>([])
        let InputText = ref('')
        let waitText = ref('')
        let TotalMessages = ref<GptMessage[]>([])
        let Loading = ref(false)
        let Waiting = ref(false)

        const scrollAreaRef = ref()

        let scrollSize = -1
        let scrollPos = 0
        let bottom = true

        const meImg = './imgs/me.jpg'
        const aiImg = './imgs/ai.png'

        DisplayMessages.value.push({
            sent: false,
            text: "你好，我是OpenAI小助手，基于gpt-3.5-turbo模型，采用ServerLess部署。<br>使用过程中有任何问题可联系：zpzhou.ok@gmail.com"
        })

        async function StreamChat() {
            if (InputText.value == "") {
                return
            }

            DisplayMessages.value.push({
                sent: true,
                text: InputText.value
            })

            TotalMessages.value.push({
                role: "user",
                content: InputText.value
            })
            InputText.value = ""
            waitText.value = ""

            // 流式聊天
            Loading.value = true
            const response = await fetch('/api/streamchat', {
                method: 'POST',
                headers: {
                    'content-type': 'application/json'
                },
                body: JSON.stringify({
                    "model": "gpt-3.5-turbo",
                    "messages": TotalMessages.value,
                    "stream": true,
                    "temperature": 0.7,
                })
            })

            const reader = response.body!.getReader()
            const decoder = new TextDecoder('utf-8')

            while (true) {
                const {value, done} = await reader.read()
                Loading.value = false
                Waiting.value = true

                if (value) {
                    let text = decoder.decode(value)
                    waitText.value = waitText.value + text
                }

                if (done) {
                    Waiting.value = false
                    TotalMessages.value?.push({
                        role: "assistant",
                        content: waitText.value
                    })
                    DisplayMessages.value.push({
                        sent: false,
                        text: waitText.value
                    })
                    break
                }
            }
        }

        function autoScroll() {
            const scroller = scrollAreaRef.value.getScroll()
            if (scroller.verticalPosition < scrollPos) {
                bottom = false
            }

            if (scroller.verticalPercentage == 1) {
                bottom = true
            }

            if (bottom) {
                scrollAreaRef.value.setScrollPercentage('vertical', 1)
            }

            scrollSize = scroller.verticalSize
            scrollPos = scroller.verticalPosition
        }

        function handleEnter(e: any) {
            if (e.ctrlKey) {
                StreamChat()
            }
        }

        return {
            handleEnter,
            StreamChat,
            autoScroll,
            scrollAreaRef,
            InputText,
            waitText,
            DisplayMessages,
            Loading,
            Waiting,
            meImg,
            aiImg
        }
    },
});
</script>
