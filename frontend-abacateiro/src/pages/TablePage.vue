<template>
  <q-page padding>
    <q-btn
      label="Create User"
      color="primary"
      @click="showModal = true"
      icon="person_add"
    />
    <!-- Create User Modal -->
    <q-dialog v-model="showModal" :backdrop-filter="'blur(4px)'">
      <q-card>
        <q-card-section>
          <div class="text-h6">Create User</div>
        </q-card-section>

        <q-card-section>
          <q-form @submit.prevent="submitForm">
            <q-input class="q-pa-sm" v-model="form.name" label="Name" filled />
            <q-input
              class="q-pa-sm"
              v-model="form.email"
              label="Email"
              type="email"
              filled
            />
            <q-input
              class="q-pa-sm"
              v-model="form.document"
              label="RG"
              type="GovernmentID"
              filled
            />
            <q-input
              v-model="form.password"
              class="q-pa-sm"
              label="Password"
              type="password"
              filled
            />
            <q-input
              v-model="form.confirmPassword"
              class="q-pa-sm"
              label="Confirm Password"
              type="password"
              filled
            />
            <q-card-actions>
              <q-btn
                flat
                label="Cancel"
                color="secondary"
                @click="showModal = false"
              />
              <q-btn label="Submit" color="primary" type="submit" />
            </q-card-actions>
          </q-form>
        </q-card-section>
      </q-card>
    </q-dialog>

        <!-- Edit User Modal -->
    <q-dialog v-model="editModal" :backdrop-filter="'blur(4px)'" >
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Edit User</div>
        </q-card-section>

        <q-card-section>
          <q-input v-model="editUser.user_name" label="Name" />
          <q-input v-model="editUser.user_email" label="Email" />
          <q-input v-model="editUser.user_document" label="RG" />
        </q-card-section>

        <q-card-actions align="right">
          <q-btn flat label="Cancel" color="primary" v-close-popup />
          <q-btn flat label="Save" color="primary" @click="saveEdit" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-table title="Usuários" :rows="rows" :columns="columns" row-key="name">
      <template v-slot:body-cell-actions="props">
        <q-td :props="props">
          <q-btn icon="mode_edit" @click="onEdit(props.row)"></q-btn>
          <q-btn icon="delete" @click="onDelete(props.row)"></q-btn>
        </q-td>
      </template>
    </q-table>

  </q-page>
</template>

<script setup>
import UserTable from "src/components/table/UserTable.vue";
import { ref, onMounted } from "vue";
import axios from "axios";
import { useRouter } from "vue-router";

const router = useRouter();
const showModal = ref(false);
const editModal = ref(false);
const rows = ref([]);
const error = ref(null);

const form = ref(createEmptyForm());
const editUser = ref(createEmptyForm());

const columns = [
  { name: 'userName', align: 'left', label: 'Name', field: row => row.user_name, sortable: true },
  { name: 'userEmail', align: 'left', label: 'Email', field: row => row.user_email, sortable: true },
  { name: 'userGovtID', align: 'left', label: 'RG', field: row => row.user_document, sortable: true },
  { name: 'actions', label: 'Action' }
];

onMounted(fetchUsers);

function createEmptyForm() {
  return { name: "", email: "", document: "", password: "", confirmPassword: "" };
}

async function submitForm() {
  if (form.value.password !== form.value.confirmPassword) {
    alert("Passwords do not match!");
    return;
  }

  try {
    await axios.post("http://localhost:8080/users", {
      user_name: form.value.name,
      user_email: form.value.email,
      user_password: form.value.password,
      user_document: form.value.document,
    });
    alert("User created successfully!");
    router.go();
  } catch (error) {
    console.error("Error creating user:", error);
    alert("Failed to create user.");
  }

  showModal.value = false;
}

async function fetchUsers() {
  try {
    const response = await axios.get('http://localhost:8080/users');
    rows.value = response.data;
  } catch (err) {
    error.value = 'Failed to load user roles';
    console.error(err);
  }
}

function onEdit(row) {
  editUser.value = { ...row };
  editModal.value = true;
}

async function saveEdit() {
  try {
    await axios.put(`http://localhost:8080/users/${editUser.value.id}`, {
      user_name: editUser.value.user_name,
      user_email: editUser.value.user_email,
      user_document: editUser.value.user_document
    });
    await fetchUsers();
    editModal.value = false;
  } catch (err) {
    console.error('Failed to update user:', err);
  }
}

async function onDelete(row) {
  try {
    await axios.delete(`http://localhost:8080/users/${row.id}`);
    await fetchUsers();
  } catch (err) {
    console.error('Failed to delete user:', err);
  }
}

defineOptions({ name: "TablePage" });
</script>
