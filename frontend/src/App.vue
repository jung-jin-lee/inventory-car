<template>
  <v-app>
    <v-navigation-drawer permanent app>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title center class="title">
            <v-img class="mx-auto" width="64" height="64" contain :aspect-ratio="16/9" src="./assets/logo.png"></v-img>
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>

      <v-divider></v-divider>

      <v-list
        dense
        nav
      >
        <v-list-item
          v-for="item in items"
          :key="item.title"
          link
        >
          <v-list-item-icon>
            <v-icon>{{ item.icon }}</v-icon>
          </v-list-item-icon>

          <v-list-item-content>
            <v-list-item-title>{{ item.title }}</v-list-item-title>
          </v-list-item-content>
        </v-list-item>
      </v-list>
    </v-navigation-drawer>

    <v-content>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  components: {
  },

  data: () => ({
    items: [
        { title: 'Inventory', },
        { title: 'Commission', },
        { title: 'Manage Market', },
        { title: 'Manage Customer', },
        { title: 'Report Setting', },
        { title: 'Sign Out', },
      ],
      right: null,
  }),
};
</script>

<style scoped>
.mx-auto {
  margin: 0 auto;
}
</style>