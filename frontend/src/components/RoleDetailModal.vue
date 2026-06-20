<script setup>
import { computed } from "vue";
import ModalDialog from "./ModalDialog.vue";
import { ShieldCheck, ShieldAlert } from "@lucide/vue";

const props = defineProps({
    visible: { type: Boolean, required: true },
    role: { type: Object, default: null },
    groupedPermissions: { type: Object, default: () => ({}) }
});

const emit = defineEmits(["close"]);

function handleClose() {
    emit("close");
}

const translateModule = (moduleName) => {
    const map = {
        'employee': 'Nhân viên',
        'user': 'Tài khoản',
        'department': 'Phòng ban',
        'role': 'Vai trò & Quyền hạn'
    };
    return map[moduleName] || moduleName;
};

// Lọc các quyền mà vai trò này đang sở hữu để hiển thị theo nhóm
const roleGroupedPermissions = computed(() => {
    if (!props.role || !props.role.permissions || props.role.permissions.length === 0) return {};
    
    // Nếu là admin, tất cả quyền
    if (props.role.name.toLowerCase() === 'admin') {
        return { 'admin': [{ code: 'All', description: 'Toàn quyền truy cập hệ thống' }] };
    }

    const groups = {};
    props.role.permissions.forEach(permCode => {
        const moduleName = permCode.split('.')[0];
        if (!groups[moduleName]) {
            groups[moduleName] = [];
        }
        groups[moduleName].push(permCode);
    });
    return groups;
});

const isSystemRole = (name) => {
    if (!name) return false;
    const l = name.toLowerCase();
    return l === 'admin' || l === 'hr' || l === 'employee';
};
</script>

<template>
    <ModalDialog
        :visible="visible"
        title="Chi tiết vai trò"
        size="lg"
        @close="handleClose"
    >
        <div class="detail-body" v-if="role">
            <!-- Header Info -->
            <div class="detail-header">
                <div class="detail-icon" :class="{'detail-icon--system': isSystemRole(role.name)}">
                    <ShieldCheck class="detail-icon-svg" />
                </div>
                <div class="detail-title-info">
                    <h3 class="detail-name">
                        {{ role.name }}
                    </h3>
                    <span v-if="isSystemRole(role.name)" class="role-system-badge">Vai trò hệ thống</span>
                </div>
            </div>

            <!-- Description -->
            <div class="detail-section">
                <h4 class="section-title">Mô tả chi tiết</h4>
                <div class="description-box">
                    {{ role.description || "Không có mô tả chi tiết cho vai trò này." }}
                </div>
            </div>

            <!-- Permissions -->
            <div class="detail-section">
                <h4 class="section-title">Quyền hạn được cấp ({{ role.name.toLowerCase() === 'admin' ? 'Toàn quyền' : (role.permissions ? role.permissions.length : 0) }})</h4>
                
                <div class="permissions-container" v-if="Object.keys(roleGroupedPermissions).length > 0">
                    <div v-for="(perms, moduleName) in roleGroupedPermissions" :key="moduleName" class="permission-group">
                        <h5 class="permission-group-title">
                            {{ moduleName === 'admin' ? 'Hệ thống' : translateModule(moduleName) }}
                        </h5>
                        <div class="permission-grid">
                            <div v-for="perm in perms" :key="perm.code || perm" class="permission-item">
                                <span class="permission-code">{{ perm.code || perm }}</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="empty-permissions" v-else>
                    <ShieldAlert class="empty-icon" />
                    <p>Vai trò này chưa được phân bổ quyền hạn nào.</p>
                </div>
            </div>
        </div>

        <template #footer>
            <button class="btn btn-secondary" @click="handleClose">Đóng</button>
        </template>
    </ModalDialog>
</template>

<style scoped>
.detail-icon--system {
    background: rgba(225, 29, 72, 0.08);
    color: var(--danger-color);
}

.role-system-badge {
    display: inline-block;
    padding: 2px 8px;
    background: rgba(225, 29, 72, 0.1);
    color: var(--danger-color);
    border-radius: var(--radius-sm);
    font-size: 11px;
    font-weight: var(--fw-medium);
    width: fit-content;
}



.description-box {
    background: var(--bg-main);
    padding: var(--space-3);
    border-radius: var(--radius-md);
    color: var(--text-main);
    line-height: 1.6;
    font-size: var(--fs-sm);
}

.permissions-container {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    max-height: 50vh;
    overflow-y: auto;
    padding-right: var(--space-2);
}

.permission-group {
    background: var(--bg-lighter);
    padding: var(--space-3);
    border-radius: var(--radius-md);
    border: 1px solid var(--border-color);
}

.permission-group-title {
    font-size: var(--fs-sm);
    font-weight: var(--fw-semibold);
    color: var(--text-main);
    margin: 0 0 var(--space-2) 0;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border-color);
    text-transform: capitalize;
}

.permission-grid {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
}

.permission-item {
    background: #ffffff;
    border: 1px solid var(--border-color);
    border-radius: var(--radius-sm);
    padding: 6px 12px;
    font-size: 12px;
    font-weight: var(--fw-medium);
    color: var(--text-main);
    box-shadow: 0 1px 2px rgba(0,0,0,0.02);
}

.permission-code {
    font-family: monospace;
}

.empty-permissions {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: var(--space-4);
    background: var(--bg-main);
    border-radius: var(--radius-md);
    color: var(--text-muted);
}

.empty-icon {
    width: 48px;
    height: 48px;
    color: var(--border-hover);
    margin-bottom: var(--space-2);
}

.permissions-container::-webkit-scrollbar {
    width: 6px;
}
.permissions-container::-webkit-scrollbar-track {
    background: transparent;
}
.permissions-container::-webkit-scrollbar-thumb {
    background-color: var(--border-hover);
    border-radius: 20px;
}
</style>
