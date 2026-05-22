<template>
   <form @submit.prevent="handleSubmit" class="department-form">
      <div class="form-grid">
         <div class="form-group">
            <label class="form-label">Tên phòng ban <span class="required">*</span></label>
            <input
               v-model="formData.name"
               type="text"
               class="form-control"
               placeholder="Nhập tên phòng ban..."
               required
            />
         </div>
         <div class="form-group">
            <label class="form-label">Mã phòng ban <span class="required">*</span></label>
            <input
               v-model="formData.code"
               type="text"
               class="form-control"
               placeholder="Nhập mã phòng ban..."
               required
            />
         </div>
         <div class="form-group form-group--full">
            <label class="form-label">Mô tả</label>
            <textarea
               v-model="formData.description"
               class="form-control"
               placeholder="Mô tả ngắn về phòng ban"
               rows="4"
            ></textarea>
         </div>
      </div>

      <div class="form-actions">
         <button type="button" class="btn btn--secondary" @click="$emit('cancel')">Hủy bỏ</button>
         <button type="submit" class="btn btn--primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ isEdit ? "Lưu thay đổi" : "Thêm phòng ban" }}
         </button>
      </div>
   </form>
</template>

<script setup>
   import { ref, watch } from "vue";

   const props = defineProps({
      initialData: { type: Object, default: () => ({}) },
      isEdit: Boolean,
      loading: Boolean,
   });

   const emit = defineEmits(["submit", "cancel"]);

   const formData = ref({ name: "", code: "", description: "" });

   // Sync form khi initialData thay đổi (immediate để load dữ liệu ngay khi mount)
   watch(
      () => props.initialData,
      (value) => {
         const data = value ?? {};
         formData.value = {
            name: data.name ?? "",
            code: data.code ?? "",
            description: data.description ?? "",
         };
      },
      { immediate: true, deep: true },
   );

   function handleSubmit() {
      emit("submit", formData.value);
   }
</script>

<style scoped>
   .department-form {
      display: flex;
      flex-direction: column;
      gap: var(--space-4);
   }

   .form-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: var(--space-3);
   }

   .form-group {
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
   }

   .form-group--full {
      grid-column: 1 / -1;
   }

   .form-label {
      font-size: var(--fs-sm);
      font-weight: var(--fw-semibold);
   }

   .form-control {
      font-size: var(--fs-base);
   }

   .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: var(--space-2);
      margin-top: var(--space-2);
      padding-top: var(--space-3);
      border-top: 1px solid var(--border-color);
   }

   .spinner {
      width: 16px;
      height: 16px;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-right-color: transparent;
      border-radius: 50%;
      animation: spin 0.8s linear infinite;
      display: inline-block;
   }

   @keyframes spin {
      to {
         transform: rotate(360deg);
      }
   }

   @media (max-width: 640px) {
      .form-grid {
         grid-template-columns: 1fr;
      }
   }
</style>
