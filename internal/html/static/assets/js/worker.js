import { hexToBuf } from 'bigint-conversion'

(async () => {
  importScripts("/assets/js/wasm_exec.js")
  const go = new Go()
  const { instance } = await WebAssembly.instantiateStreaming(
    fetch('/assets/wasm/srp.wasm', {
      integrity: 'sha384-GvKYPyi99AXVg0Lhgf0Wmz+/5capMDB7tBjqYm1+bQT0hwsVltHf113N8guk7wdm'
    }),
    go.importObject
  )
  await go.run(instance)
}
)()

self.onmessage = async function(e) {
  const [rid, action, args] = e.data
  const { identifier, password } = args
  try {
    switch (action.toLowerCase()) {
      case 'register':
        const { salt, verifier } = await getRegistrationInfo(identifier, password)
        const response = await fetch("/api/auth/register", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            "identifier": identifier,
            "salt": salt,
            "verifier": verifier,
          }),
        })

        const data = await response.json()
        if (response.status !== 200) {
          throw new Error(data.error)
        }

        self.postMessage([rid, true, true])
      case 'login':
        const identityResp = await fetch("/api/auth/identify", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            "identifier": identifier,
          }),
        })

        const idData = await identityResp.json()
        if (identityResp.status !== 200) {
          throw new Error(idData.error)
        }

        const saltU8 = new Uint8Array(hexToBuf(idData['salt']))
        const B = new Uint8Array(hexToBuf(idData['B']))
        const A = await setupClient(identifier, password, saltU8, B)
        const proof = await getClientProof()

        const loginResp = await fetch("/api/auth/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            "identifier": identifier,
            "A": A,
            "proof": proof,
          }),
        })

        const loginData = await loginResp.json()
        if (loginResp.status !== 200) {
          throw new Error(loginData.error)
        }

        const serverProof = new Uint8Array(hexToBuf(loginData['proof']))

        const valid = await verifyServerProof(serverProof)
        if (!valid) {
          throw new Error('invalid server proof')
        }

        const key = await getKey()
        self.postMessage([rid, true, key])
      default:
        self.postMessage([rid, false, new Error('invalid action')])
    }
  } catch (err) {
    self.postMessage([rid, false, err])
  }
}

self.postMessage('ready')
