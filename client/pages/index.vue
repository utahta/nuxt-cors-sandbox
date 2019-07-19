<template>
  <div class="container">
    <div>
      <p> SessionID: {{ session_id }} </p>
      <button @click="callAPI">Click</button>
    </div>
  </div>
</template>

<script>
import Logo from '~/components/Logo.vue'
import axios from 'axios'

export default {
  components: {
    Logo
  },
  computed: {
    session_id() {
      return this.$store.state.session_id
    }
  },
  methods: {
    callAPI() {
      let ctx = this
      let c = axios.create({
        baseURL: 'https://localhost:8080',
        withCredentials: true,
      })
      c.post('/api/hello', {}, {
        'headers': {'Authorization': 'aaaaa'}
      }).then(function(resp) {
          ctx.$store.commit('setSessionID', ctx.$store.state.session_id +  ' + ' + resp.data)
        })
    }
  }
}
</script>

<style>
.container {
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.title {
  font-family: 'Quicksand', 'Source Sans Pro', -apple-system, BlinkMacSystemFont,
    'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}
</style>
