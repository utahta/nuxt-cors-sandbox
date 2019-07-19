import axios from 'axios'
import https from 'https'

export const state = () => {
  return {
    session_id: null
  }
}
export const mutations = {
  setSessionID (state, sid) {
    state.session_id = sid
  }
}
export const actions = {
  nuxtServerInit() {
    console.log('nuxt server init')
    let c = axios.create({
      baseURL: 'https://localhost:8080',
      withCredentials: true,
      httpsAgent: new https.Agent({
        rejectUnauthorized: false // for local
      })
    })
    c.post('/api/hello', {}, {
      'headers': {'Cookie': 'session_id=inject;'}
    }).then(function(resp) {
      console.log('call api/hello')
    })
  },
  setSession({ commit }) {
    console.log('set session')
    let sid = this.$cookies.get('session_id')
    if (sid) {
      console.log('commit session')
      commit('setSessionID', sid)
    }
  },
}
export const getters = {
  isLogin: state => {
    console.log('isLogin')
    return state.session_id !== null
  }
}
