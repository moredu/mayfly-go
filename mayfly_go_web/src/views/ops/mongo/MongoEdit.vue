<template>
    <div>
        <el-dialog :title="title" v-model="dialogVisible" :before-close="cancel" :close-on-click-modal="false" width="38%" :destroy-on-close="true">
            <el-form :model="form" ref="mongoForm" :rules="rules" label-width="85px">
                <el-tabs v-model="tabActiveName">
                    <el-tab-pane label="基础信息" name="basic">
                        <el-form-item ref="tagSelectRef" prop="tagCodePaths" label="标签" required>
                            <tag-tree-select
                                @change-tag="
                                    (tagCodePaths) => {
                                        form.tagCodePaths = tagCodePaths;
                                        tagSelectRef.validate();
                                    }
                                "
                                multiple
                                :select-tags="form.tagCodePaths"
                                style="width: 100%"
                            />
                        </el-form-item>
                        <el-form-item prop="code" label="编号" required>
                            <el-input
                                :disabled="form.id"
                                v-model.trim="form.code"
                                placeholder="请输入编号 (大小写字母、数字、_-.:), 不可修改"
                                auto-complete="off"
                            ></el-input>
                        </el-form-item>
                        <el-form-item prop="name" label="名称" required>
                            <el-input v-model.trim="form.name" placeholder="请输入名称" auto-complete="off"></el-input>
                        </el-form-item>
                        <el-form-item prop="uri" label="uri" required>
                            <el-input
                                type="textarea"
                                :rows="2"
                                v-model.trim="form.uri"
                                placeholder="形如 mongodb://username:password@host1:port1"
                                auto-complete="off"
                            ></el-input>
                        </el-form-item>
                    </el-tab-pane>

                    <el-tab-pane label="其他配置" name="other">
                        <el-form-item prop="sshTunnelMachineId" label="SSH隧道">
                            <ssh-tunnel-select v-model="form.sshTunnelMachineId" />
                        </el-form-item>
                    </el-tab-pane>
                </el-tabs>
            </el-form>

            <template #footer>
                <div class="dialog-footer">
                    <el-button @click="testConn" :loading="testConnBtnLoading" type="success">测试连接</el-button>
                    <el-button @click="cancel()">取 消</el-button>
                    <el-button type="primary" :loading="saveBtnLoading" @click="btnOk">确 定</el-button>
                </div>
            </template>
        </el-dialog>
    </div>
</template>

<script lang="ts" setup>
import { toRefs, reactive, ref, watchEffect } from 'vue';
import { mongoApi } from './api';
import { ElMessage } from 'element-plus';
import TagTreeSelect from '../component/TagTreeSelect.vue';
import SshTunnelSelect from '../component/SshTunnelSelect.vue';
import { ResourceCodePattern } from '@/common/pattern';

const props = defineProps({
    visible: {
        type: Boolean,
    },
    mongo: {
        type: [Boolean, Object],
    },
    title: {
        type: String,
    },
});

//定义事件
const emit = defineEmits(['update:visible', 'cancel', 'val-change']);

const rules = {
    tagCodePaths: [
        {
            required: true,
            message: '请选择标签',
            trigger: ['change', 'blur'],
        },
    ],
    code: [
        {
            required: true,
            message: '请输入编码',
            trigger: ['change', 'blur'],
        },
        {
            pattern: ResourceCodePattern.pattern,
            message: ResourceCodePattern.message,
            trigger: ['blur'],
        },
    ],
    name: [
        {
            required: true,
            message: '请输入名称',
            trigger: ['change', 'blur'],
        },
    ],
    uri: [
        {
            required: true,
            message: '请输入mongo uri',
            trigger: ['change', 'blur'],
        },
    ],
};

const mongoForm: any = ref(null);
const tagSelectRef: any = ref(null);

const state = reactive({
    dialogVisible: false,
    tabActiveName: 'basic',
    form: {
        id: null,
        code: '',
        name: null,
        uri: null,
        sshTunnelMachineId: null as any,
        tagCodePaths: [],
    },
    submitForm: {},
});

const { dialogVisible, tabActiveName, form, submitForm } = toRefs(state);

const { isFetching: testConnBtnLoading, execute: testConnExec } = mongoApi.testConn.useApi(submitForm);
const { isFetching: saveBtnLoading, execute: saveMongoExec } = mongoApi.saveMongo.useApi(submitForm);

watchEffect(() => {
    state.dialogVisible = props.visible;
    if (!state.dialogVisible) {
        return;
    }
    state.tabActiveName = 'basic';
    const mongo: any = props.mongo;
    if (mongo) {
        state.form = { ...mongo };
        state.form.tagCodePaths = mongo.tags.map((t: any) => t.codePath);
    } else {
        state.form = { db: 0 } as any;
    }
});

const getReqForm = () => {
    const reqForm = { ...state.form };
    if (!state.form.sshTunnelMachineId || state.form.sshTunnelMachineId <= 0) {
        reqForm.sshTunnelMachineId = -1;
    }
    return reqForm;
};

const testConn = async () => {
    mongoForm.value.validate(async (valid: boolean) => {
        if (!valid) {
            ElMessage.error('请正确填写信息');
            return false;
        }

        state.submitForm = getReqForm();
        await testConnExec();
        ElMessage.success('连接成功');
    });
};

const btnOk = async () => {
    mongoForm.value.validate(async (valid: boolean) => {
        if (!valid) {
            ElMessage.error('请正确填写信息');
            return false;
        }
        state.submitForm = getReqForm();
        await saveMongoExec();
        ElMessage.success('保存成功');
        emit('val-change', state.form);
        cancel();
    });
};

const cancel = () => {
    emit('update:visible', false);
    emit('cancel');
};
</script>
<style lang="scss"></style>
