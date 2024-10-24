<template>
  <q-page class="q-pa-md">
    <q-btn
      label="Upload File"
      color="primary"
      @click="openModal"
      class="q-mb-md"
    />

    <q-dialog v-model="modalOpen">
      <q-card style="min-width: 350px">
        <q-card-section>
          <div class="text-h6">Upload File</div>
        </q-card-section>

        <q-card-section class="q-pt-none">
          <q-file
            v-model="file"
            label="Choose file"
            filled
            bottom-slots
            counter
            max-file-size="10485760"
            max-files="1"
            accept=".docx"
            @rejected="onRejected"
          >
            <template v-slot:prepend>
              <q-icon name="attach_file" />
            </template>
          </q-file>
        </q-card-section>

        <q-card-actions align="right" class="text-primary">
          <q-btn flat label="Cancel" v-close-popup />
          <q-btn flat label="Upload" @click="uploadFile" v-close-popup />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-table
      title="Work Reports"
      :rows="workReports"
      :columns="columns"
      row-key="id"
      :loading="loading"
    />
  </q-page>
</template>

<script setup lang="js">
import axios from "axios";
import { ref, onMounted } from "vue";

const modalOpen = ref(false);
const file = ref(null);

const openModal = () => {
  modalOpen.value = true;
};

const workReports = ref([]);
const loading = ref(false);

const columns = [
  { name: "id", label: "ID", field: "id", sortable: true },
  { name: "fileName", label: "File Name", field: "fileName", sortable: true },
  {
    name: "uploadDate",
    label: "Upload Date",
    field: "uploadDate",
    sortable: true,
  },
  // Add more columns as needed based on your backend response
];

const fetchWorkReports = async () => {
  loading.value = true;
  try {
    const response = await axios.get("http://localhost:8080/work-reports");
    workReports.value = response.data;
  } catch (error) {
    console.error("Error fetching work reports:", error);
    // You might want to add some user feedback here, like an error message
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchWorkReports();
});

const uploadFile = async () => {
  if (!file.value) {
    console.error("No file selected");
    return;
  }

  const data = new FormData();
  data.append("report", file.value);

  // const jsonData = { report: file.value };
  // formData.append("file", file.value);
  console.log(data);

  try {
    const response = await axios.post(
      "http://localhost:8080/work-reports/" + "teste",
      data,
      {
        headers: {
          "Content-Type": "multipart/form-data",
        },
      },
    );
    console.log("File uploaded successfully:", response.data);
    // Refresh the table after successful upload
    await fetchWorkReports();
    // You might want to add some user feedback here, like a success message
  } catch (error) {
    console.error("Error uploading file:", error);
    // You might want to add some user feedback here, like an error message
  }
};

const onRejected = (files) => {
  console.log(files);
};
</script>

<style scoped></style>
