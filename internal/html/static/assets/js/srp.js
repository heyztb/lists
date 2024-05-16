class SRPClient {
  rid = 0
  promises = new Map()
  worker = new Worker('/assets/js/worker.min.js')
  ready = new Promise((resolve) => {
    this.worker.addEventListener('message', (x) => x.data === 'ready' && resolve(true), { once: true })
  })

  constructor() {
    this.ready.then(() =>
      this.worker.onmessage = (e) => {
        const [rid, suc, res] = e.data;
        const [resolve, reject] = this.promises.get(rid);
        this.promises.delete(rid);
        (suc ? resolve : reject)(res)
      }
    )
  }

  async register(identifier, password) {
    return new Promise((resolve, reject) => {
      this.promises.set(this.rid, [resolve, reject])
      this.worker.postMessage([this.rid, 'register', { identifier, password }])
      this.rid++
    })
  }

  async login(identifier, password) {
    return new Promise((resolve, reject) => {
      this.promises.set(this.rid, [resolve, reject])
      this.worker.postMessage([this.rid, 'login', { identifier, password }])
      this.rid++
    })
  }

  terminate() {
    this.worker.terminate()
  }
}

export { SRPClient }