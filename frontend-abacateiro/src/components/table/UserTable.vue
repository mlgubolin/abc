<template>
  <div class="q-pa-md">
    <!-- <q-responsive> -->
    <q-table title="Usuários" :rows="rows" :columns="columns" row-key="name">
      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn icon="mode_edit" @click="onEdit(props.row)"></q-btn>
          <q-btn icon="delete" @click="onDelete(props.row)"></q-btn>
        </q-td>
      </template>
    </q-table>

    <!-- </q-responsive>  -->
  </div>
</template>

<script setup lang="js">
import { ref, onMounted } from 'vue';
import axios from 'axios';

// Call fetchUserRoles when the component is mounted

</script>
