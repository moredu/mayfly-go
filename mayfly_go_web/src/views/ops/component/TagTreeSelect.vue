<template>
    <div>
        <el-tree-select
            v-bind="$attrs"
            v-model="state.selectTags"
            @change="changeTag"
            :data="tags"
            placeholder="请选择关联标签"
            :render-after-expand="true"
            :default-expanded-keys="[state.selectTags]"
            show-checkbox
            node-key="codePath"
            :check-strictly="props.checkStrictly"
            :props="{
                value: 'codePath',
                label: 'codePath',
                children: 'children',
            }"
        >
            <template #default="{ data }">
                <span class="custom-tree-node">
                    <span style="font-size: 13px">
                        {{ data.code }}
                        <span style="color: #3c8dbc">【</span>
                        {{ data.name }}
                        <span style="color: #3c8dbc">】</span>
                        <el-tag v-if="data.children !== null" size="small">{{ data.children.length }}</el-tag>
                    </span>
                </span>
            </template>
        </el-tree-select>
    </div>
</template>

<script lang="ts" setup>
import { toRefs, reactive, onMounted } from 'vue';
import { tagApi } from '../tag/api';
import { TagResourceTypeEnum } from '@/common/commonEnum';

//定义事件
const emit = defineEmits(['update:modelValue', 'changeTag', 'input']);

const props = defineProps({
    selectTags: {
        type: [Array<any>],
    },
    tagType: {
        type: Number,
        default: TagResourceTypeEnum.Tag.value,
    },
    checkStrictly: {
        type: Boolean,
        default: false,
    },
});

const state = reactive({
    tags: [],
    // 单选则为codePath，多选为codePath数组
    selectTags: [] as any,
});

const { tags } = toRefs(state);

onMounted(async () => {
    state.selectTags = props.selectTags;
    state.tags = await tagApi.getTagTrees.request({ type: props.tagType });
});

const changeTag = () => {
    emit('changeTag', state.selectTags);
};
</script>
<style lang="scss"></style>
