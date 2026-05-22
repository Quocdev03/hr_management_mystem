<template>
   <form @submit.prevent="handleSubmit" class="employee-form">
      <div class="form-grid">
         <!-- Họ và tên -->
         <div class="form-group">
            <label class="form-label">Họ <span class="required">*</span></label>
            <input
               v-model="formData.first_name"
               type="text"
               class="form-control"
               placeholder="Nhập họ..."
               required
            />
         </div>
         <div class="form-group">
            <label class="form-label">Tên <span class="required">*</span></label>
            <input
               v-model="formData.last_name"
               type="text"
               class="form-control"
               placeholder="Nhập tên..."
               required
            />
         </div>

         <!-- Liên hệ -->
         <div class="form-group" :class="{ 'form-group--disabled': isEdit }">
            <label class="form-label">Email <span class="required" v-if="!isEdit">*</span></label>
            <input
               v-model="formData.email"
               type="email"
               class="form-control"
               placeholder="example@company.com"
               :required="!isEdit"
               :disabled="isEdit"
            />
            <small v-if="isEdit" class="form-hint">Không thể thay đổi email</small>
         </div>
         <div class="form-group">
            <label class="form-label">Số điện thoại <span class="required">*</span></label>
            <input
               v-model="formData.phone"
               type="text"
               class="form-control"
               placeholder="0xxxxxxxxx"
               required
            />
         </div>

         <!-- Công việc -->
         <div class="form-group">
            <label class="form-label">Phòng ban <span class="required">*</span></label>
            <select v-model="formData.department_id" class="form-control" required>
               <option value="" disabled>Chọn phòng ban</option>
               <option v-for="dept in departments" :key="dept.id" :value="dept.id">
                  {{ dept.name }}
               </option>
            </select>
         </div>
         
         <div class="form-group">
            <label class="form-label">Chức vụ</label>
            <input
               v-model="formData.position"
               type="text"
               class="form-control"
               placeholder="Ví dụ: Backend Developer"
            />
         </div>

         <!-- Lương và Ngày vào làm -->
         <div class="form-group">
            <label class="form-label">Mức lương (VNĐ)</label>
            <input
               v-model.number="formData.salary"
               type="number"
               class="form-control"
               placeholder="0"
            />
         </div>
         <div class="form-group" :class="{ 'form-group--disabled': isEdit }">
            <label class="form-label">Ngày vào làm</label>
            <input
               v-model="formData.join_date"
               type="date"
               class="form-control"
               :disabled="isEdit"
            />
         </div>
         <!-- Trạng thái và Tài khoản -->
         <div class="form-group">
            <label class="form-label">Trạng thái</label>
            <select v-model="formData.status" class="form-control">
               <option value="active">Đang làm việc</option>
               <option value="inactive">Đã nghỉ việc</option>
            </select>
         </div>
      </div>

      <div class="form-actions">
         <button type="button" class="btn btn--secondary" @click="$emit('cancel')">Hủy bỏ</button>
         <button type="submit" class="btn btn--primary" :disabled="loading">
            <span v-if="loading" class="spinner"></span>
            {{ isEdit ? "Lưu thay đổi" : "Thêm nhân viên" }}
         </button>
      </div>
   </form>
</template>

<script setup>
   // ============================================================
   // 1. IMPORTS
   // ============================================================
   import { ref, watch, onMounted } from "vue";
   import { useDepartmentStore } from "@/store/department";
   import { useToast } from "vue-toastification";

   // ============================================================
   // 2. PROPS & EMITS
   // ============================================================
   const props = defineProps({
      initialData: { type: Object, default: () => ({}) },
      isEdit: Boolean,
      loading: Boolean,
   });

   const emit = defineEmits(["submit", "cancel"]);

   // ============================================================
   // 3. STATE (reactive data)
   // ============================================================
   const toast = useToast();
   const deptStore = useDepartmentStore();

   const formData = ref(buildFormData(props.initialData));

   // dùng thẳng state từ store, không cần local copy
   const departments = deptStore.departments;

   // ============================================================
   // 4. HELPERS
   // ============================================================

   /** Tạo object formData từ data ban đầu, đảm bảo luôn có đủ các field */
   function buildFormData(data) {
      const d = data ?? {}; // null hoặc undefined đều fallback về {}
      return {
         first_name: d.first_name ?? "",
         last_name: d.last_name ?? "",
         email: d.email ?? "",
         phone: d.phone ?? "",
         department_id: d.department_id ?? "",
         position: d.position ?? "",
         salary: d.salary ?? null,
         join_date: d.join_date ?? "",
         status: d.status ?? "active",
      };
   }

   /** So sánh từng field và chỉ lấy những field đã thay đổi */
   function getDirtyPayload(current, original) {
      const fields = [
         "first_name",
         "last_name",
         "phone",
         "department_id",
         "position",
         "salary",
         "status",
      ];

      const payload = {};

      for (const field of fields) {
         if (current[field] !== original[field]) {
            payload[field] = current[field];
         }
      }

      return payload;
   }

   // ============================================================
   // 5. API CALLS
   // ============================================================

   async function fetchDepts() {
      // Store đã xử lý try/catch và tự cập nhật deptStore.departments
      await deptStore.fetchDepartments();
   }

   // ============================================================
   // 6. EVENT HANDLERS
   // ============================================================

   function handleSubmit() {
      if (!props.isEdit) {
         // Thêm mới: gửi toàn bộ form
         emit("submit", formData.value);
         return;
      }

      // Cập nhật: chỉ gửi những field đã thay đổi
      const originalData = buildFormData(props.initialData);
      const payload = getDirtyPayload(formData.value, originalData);

      if (Object.keys(payload).length === 0) {
         toast.info("Không có dữ liệu thay đổi");
         return;
      }

      emit("submit", payload);
   }

   // ============================================================
   // 7. WATCHERS
   // ============================================================

   // Reset form khi props.initialData thay đổi (ví dụ: mở modal edit nhân viên khác)
   watch(
      () => props.initialData,
      (newData) => {
         formData.value = buildFormData(newData ?? {});
      },
      { deep: true },
   );

   // ============================================================
   // 8. LIFECYCLE HOOKS
   // ============================================================

   onMounted(() => {
      fetchDepts();
   });
</script>

<style scoped>
   .employee-form {
      display: flex;
      flex-direction: column;
      gap: var(--space-4);
   }

   .form-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: var(--space-3);
   }

   .form-actions {
      display: flex;
      justify-content: flex-end;
      gap: var(--space-2);
      margin-top: var(--space-2);
      padding-top: var(--space-3);
      border-top: 1px solid var(--border-color);
   }

   /* Đồng bộ font chữ với main.css */
   .form-label {
      font-size: var(--fs-sm);
      font-weight: var(--fw-semibold);
   }

   .form-control {
      font-size: var(--fs-base);
   }

   .form-group--disabled {
      opacity: 0.8;
   }

   .form-group--disabled .form-control {
      background-color: var(--bg-light);
      cursor: not-allowed;
      border-color: var(--border-color);
   }

   .form-hint {
      font-size: 11px;
      color: var(--text-light);
      margin-top: -2px;
   }

   @media (max-width: 640px) {
      .form-grid {
         grid-template-columns: 1fr;
      }
   }

   .spinner {
      width: 16px;
      height: 16px;
      border: 2px solid rgba(255, 255, 255, 0.3);
      border-radius: 50%;
      border-top-color: #fff;
      animation: spin 0.8s linear infinite;
      display: inline-block;
      margin-right: 8px;
   }

   @keyframes spin {
      to {
         transform: rotate(360deg);
      }
   }
</style>
