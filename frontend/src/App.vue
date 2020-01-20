<template>
    <v-app>
        <v-navigation-drawer permanent app sm="400" md="320" lg="800" :width="navigationMenuWidth">
            <v-list-item>
                <v-list-item-title center class="title">
                    <v-row>
                        <v-col cols="12">
                            <v-img class="mx-auto" width="64" height="64" contain :aspect-ratio="16/9" src="./assets/logo.png"></v-img>
                        </v-col>
                    </v-row>
                </v-list-item-title>
            </v-list-item>

            <v-divider></v-divider>

            <v-list
              nav
            >
                <v-list-item
                  :key="item.title"
                  link
                  v-for="item in items"
                >
                    <v-list-item-icon>
                        <v-icon>{{ item.icon }}</v-icon>
                    </v-list-item-icon>

                    <v-list-item-content>
                        <router-link :to="item.path" v-if="item.path">
                            <v-list-item-title>{{ item.title }}</v-list-item-title>
                        </router-link>
                        <v-list-item-title v-else>{{ item.title }}</v-list-item-title>
                    </v-list-item-content>
                </v-list-item>
            </v-list>
        </v-navigation-drawer>

        <v-content>
            <v-container fluid>
                <router-view></router-view>
            </v-container>
        </v-content>
    </v-app>
</template>

<script>

export default {
  name: 'App',
  computed: {
    navigationMenuWidth() {
      switch(this.$vuetify.breakpoint.name) {
        case "xs": return "240px"
        case "sm": return "240px"
        case "md": return "280px"
        case "lg": return "320px"
        case "xl": return "400px"
      }
      return "240px";
    },
  },
  data: () => ({
    items: [
        { "title": 'Inventory', "path": "/inventory/cars" },
        { "title": 'Commission', "path": "/commision" },
        { "title": 'Manage Market', "path": "/manage/market" },
        { "title": 'Manage Customer', },
        { "title": 'Report Setting', },
        { "title": 'Sign Out', },
      ],
      right: null,
  }),
};
</script>

<style scoped>
.mx-auto {
  margin: 0 auto;
}
.v-application a,
.router-link-active {
  color: inherit;
  text-decoration: none;
  outline: none;
}
</style>