<template>
  <div id="app">
    <nav :class="{ collapsed: isNavCollapsed }">
      <router-link :to="{ name: 'Home' }">
        <img :src="navLogo" alt="Logo" class="logo" />
      </router-link>
      <router-link :to="{ name: 'Maps' }">
        {{ isNavCollapsed ? "" : "Maps" }}<i class="fa-solid fa-map-pin"></i>
      </router-link>
      <router-link :to="{ name: 'Gowler' }">
        {{ isNavCollapsed ? "" : "Gowler" }}<i class="fa-solid fa-spider"></i>
      </router-link>

      <button @click="toggleNav" class="toggle-btn">
        <i :class="isNavCollapsed ? 'fa-solid fa-square-caret-right' : 'fa-solid fa-square-caret-left'"></i>
      </button>
    </nav>
    <div id="body">
      <router-view/>
    </div>
  </div>
</template>

<script>
import logoCollapsed from '@/assets/imgs/logo_25x25.png';
import logoExpanded from '@/assets/imgs/logo_120x120.png';

export default {
  data() {
    return {
      isNavCollapsed: false,
    };
  },
  methods: {
    toggleNav() {
      this.isNavCollapsed = !this.isNavCollapsed;
    },
  },
  computed: {
    navLogo() {
      return this.isNavCollapsed ? logoCollapsed : logoExpanded;
    },
  },
};
</script>

<style>
* {
  box-sizing: border-box;
}
#app {
  display: flex;
  height: 100vh;
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  height: 100vh;
}

nav {
  width: 10vw;
  height: 100%;
  display: flex;
  flex-flow: column;
  background-color: #006D77;
  padding: 50px 35px 50px 25px;
  transition: width 0.3s ease;
}

nav.collapsed {
  width: 50px;
}

nav a {
  display: flex;
  justify-content: space-between;
  color: #EDF6F9;
  text-decoration: none;
  font-size: 18px;
  margin: 10px 0;
  opacity: 1;
  transition: opacity 0.3s ease;
}

nav a:hover {
  color: #83C5BE;
}

nav a i {
  width: 18px;
  height: 18px;
  margin-left: auto;
}

nav a img {
  display: block;
  margin: 0 auto;
}

.toggle-btn {
  width: 18px;
  height: 18px;
  background: none;
  border: none;
  color: #EDF6F9;
  font-size: 18px;
  cursor: pointer;
  margin-bottom: 20px;
  margin-left: auto;
  padding: 0;
}

.toggle-btn:hover {
  color: #83C5BE;
}

#body {
  padding: 50px 25px;
  overflow-y: auto;
  flex: 1;
  width: 90vw;
}
</style>
