<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-toolbar-title>
          Let'sChat
        </q-toolbar-title>
          by zpzhou
      </q-toolbar>

    </q-header>

    <q-page-container>
      <router-view/>
    </q-page-container>
  </q-layout>

</template>

<script lang="ts">
import {defineComponent} from 'vue';


export default defineComponent({
  name: 'MainLayout',

  setup() {
    return {}
  }
});
</script>
