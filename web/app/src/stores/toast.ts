import { defineStore } from 'pinia'

type ToastTone = 'success' | 'error'

let toastId = 0

export const useToastStore = defineStore('toast', {
  state: () => ({
    id: 0,
    message: '',
    tone: 'success' as ToastTone,
    timer: 0
  }),
  actions: {
    show(message: string, tone: ToastTone = 'success') {
      if (this.timer) window.clearTimeout(this.timer)
      this.id = ++toastId
      this.message = message
      this.tone = tone
      this.timer = window.setTimeout(() => this.clear(), 2400)
    },
    clear() {
      this.message = ''
      this.timer = 0
    }
  }
})
