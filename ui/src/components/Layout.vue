<script lang="ts" setup>
import { useOsTheme, darkTheme } from 'naive-ui'
import { computed } from 'vue';
const osThemeRef = useOsTheme()
let theme = computed(() => (osThemeRef.value === 'dark' ? darkTheme : null))

</script>
    
<template>
    <div>
        <n-config-provider :theme="theme">
            <n-global-style />
            <n-space vertical size="large">
                <n-card> 当前操作系统的主题是 {{ osThemeRef }}。 </n-card>
            </n-space>
        </n-config-provider>
    </div>
</template>
    