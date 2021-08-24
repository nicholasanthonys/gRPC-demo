<template>
  <img alt="Vue logo" src="./assets/logo.png">
  <HelloWorld msg="Welcome to Your Vue.js App"/>
</template>

<script>
import HelloWorld from './components/HelloWorld.vue'

// import {ProtoGrpcType} from 'proto/greet'
import {GreetServiceClient,} from './proto/greet_grpc_web_pb'
import {GreetRequest, Greeting, } from './proto/greet_pb'
// import {Greeting} from '../proto/greet/Greeting'

export default {
  name: 'App',
  components: {
    HelloWorld
  },
  methods: {
    doUnaryRequest() {


      const client = new GreetServiceClient('http://' + window.location.hostname + ':8080', null, null)


      // make request to grpc
      const request = new GreetRequest()
      let greeting = new Greeting()
      greeting.setFirstName("Spiderman")
      greeting.setLastName("Marvel")
      request.setGreeting(greeting)
      console.log("request is")
      console.log(request)

      try {
        client.greet(request, null, (err, response) => {
          console.log("result : ")
          console.log(response.getResult())
        })


      } catch (err) {
        console.log("error : ")
        console.log(err)
      }

    },


  },
  mounted() {
    this.doUnaryRequest()
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
