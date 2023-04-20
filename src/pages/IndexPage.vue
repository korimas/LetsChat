<template>
    <q-page style="min-height: 0;">
        <!--密码认证-->
        <q-dialog v-model="AuthRequire" persistent>
            <q-card style="min-width: 350px">
                <q-card-section>
                    <div class="text-h6">请输入密码</div>
                </q-card-section>

                <q-card-section class="q-pt-none">
                    <q-input dense type="password" v-model="Password" autofocus @keyup.enter="Auth" />
                    <p class="text-grey text-right">Note: My company name.</p>
                </q-card-section>

                <q-card-actions align="right" class="text-primary">
                    <q-btn flat label="确定" @click="Auth"/>
                </q-card-actions>
            </q-card>
        </q-dialog>

        <div class="q-pa-md column no-wrap" v-if="AuthFinish">

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
            <div class="row justify-center" >
                <div class="row justify-center" style="width: 100%; max-width: 800px">
                    <q-input
                        autofocus
                        ref="inputCom"
                        square
                        class="col-10"
                        :disable="Loading || Waiting"
                        @keydown.enter="handleEnter"
                        filled autogrow bg-color="grey"
                        v-model="InputText"
                    />

                    <q-btn
                        square
                        @click="StreamChat"
                        unelevated
                        class="col-2"
                        color="secondary"
                    >
                        <div>发送</div>
                    </q-btn>
                </div>
            </div>
        </div>
    </q-page>
</template>

<script lang="ts">
import {defineComponent, ref, nextTick} from 'vue';
import api from 'src/api/request'
import { QInput } from 'quasar'

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
        const inputCom = ref(QInput)
        let DisplayMessages = ref<Message[]>([])
        let InputText = ref('')
        let waitText = ref('')
        let TotalMessages = ref<GptMessage[]>([])
        let Loading = ref(false)
        let Waiting = ref(false)
        let AuthFinish = ref(false)
        let AuthRequire = ref(true)
        let Password = ref('')

        const scrollAreaRef = ref()

        let scrollSize = -1
        let scrollPos = 0
        let needBottom = true

        const meImg = './imgs/me.jpg'
        const aiImg = './imgs/ai.png'

        DisplayMessages.value.push({
            sent: false,
            text: "Hello，我是OpenAI小助手，基于gpt-3.5-turbo模型。"
        })
        DisplayMessages.value.push({
            sent: false,
            text: "向聊天框发送信息即可与我聊天，我可以回答问题，写代码、写文章等等"
        })
        DisplayMessages.value.push({
            sent: false,
            text: "使用过程中有任何问题可联系：zpzhou.ok@gmail.com"
        })

        checkAuth()

        function checkAuth() {
            api.CheckNeedAuth().then(response => {
                AuthRequire.value = response.data.result.authRequire
                // AuthRequire.value = true

                AuthFinish.value = !AuthRequire.value
            })
        }

        function Auth() {
            api.PasswordAuth(Password.value).then(response => {
                if (response.data.success){
                    AuthFinish.value = true
                    AuthRequire.value = false
                }
            })
        }

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
            await nextTick()
            scrollBottom()

            // 流式聊天
            Loading.value = true
            const response = await fetch('/api/streamchat', {
                method: 'POST',
                headers: {
                    'content-type': 'application/json',
                    'Authorization': 'Bearer ' + Password.value
                },
                body: JSON.stringify({
                    "model": "gpt-3.5-turbo",
                    "messages": TotalMessages.value,
                    "stream": true,
                    "temperature": 0.7,
                })
            })

            scrollBottom()
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
                    await nextTick()
                    inputCom.value.focus()
                    break
                }
            }
        }

        function scrollBottom() {
            scrollAreaRef.value.setScrollPercentage('vertical', 1)
        }

        function autoScroll() {
            const scroller = scrollAreaRef.value.getScroll()
            console.log("verticalPosition: " + scroller.verticalPosition)
            console.log("verticalSize: " + scroller.verticalSize)
            console.log("verticalContainerSize: " + scroller.verticalContainerSize)
            console.log("horizontalSize: " + scroller.horizontalSize)
            console.log("horizontalContainerSize: " + scroller.horizontalContainerSize)
            console.log("verticalPercentage: " + scroller.verticalPercentage)

            console.log("---------------------------------------------------------")

            if (needBottom) {
                scrollBottom()
                needBottom = false
                return
            }

            if (scroller.verticalSize - scroller.verticalContainerSize - scroller.verticalPosition < 120) {
                scrollBottom()
            }

        }

        function handleEnter(e: any) {
            if (!e.ctrlKey) {
                const scroller = scrollAreaRef.value.getScroll()
                if (scroller.verticalSize - scroller.verticalContainerSize - scroller.verticalPosition < 120) {
                    needBottom=true
                }

                StreamChat()
            } else {
                InputText.value = InputText.value + "\n"
            }
        }

        return {
            handleEnter,
            StreamChat,
            autoScroll,
            scrollAreaRef,
            InputText,
            waitText,
            Password,
            DisplayMessages,
            Loading,
            Waiting,
            AuthFinish,
            AuthRequire,
            Auth,
            meImg,
            aiImg,
            inputCom
        }
    },
});
</script>
