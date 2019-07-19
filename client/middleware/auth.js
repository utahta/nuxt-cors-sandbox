export default function ({ store, redirect }) {
  console.log('auth')
  store.dispatch('setSession')

  if (store.getters.isLogin === false) {
    console.log('redirect')
    return redirect('/auth')
  }
}
