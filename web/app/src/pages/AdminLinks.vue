<script setup lang="ts">
import { computed, onMounted, reactive } from 'vue'
import { Plus, Save, Trash2 } from 'lucide-vue-next'
import { useI18n } from 'vue-i18n'
import type { NavLink } from '../api/types'
import { useNavigationStore } from '../stores/navigation'
import { useToastStore } from '../stores/toast'

const nav = useNavigationStore()
const toast = useToastStore()
const { t } = useI18n()
const newGroup = reactive({ title: '', sortOrder: 100, visible: true })
const newLink = reactive<Partial<NavLink>>({
  groupId: 0,
  title: '',
  url: '',
  icon: '',
  iconType: 'text',
  description: '',
  sortOrder: 100,
  visible: true,
  openInNewTab: true
})

const flatLinks = computed(() => nav.groups.flatMap(group => group.links))

onMounted(async () => {
  await nav.loadAdmin()
  if (!newLink.groupId && nav.groups[0]) newLink.groupId = nav.groups[0].id
})

async function addGroup() {
  await nav.saveGroup(newGroup)
  newGroup.title = ''
  newGroup.sortOrder += 10
  toast.show(t('addSuccess'))
}

async function addLink() {
  await nav.saveLink(newLink)
  Object.assign(newLink, {
    groupId: nav.groups[0]?.id || 0,
    title: '',
    url: '',
    icon: '',
    iconType: 'text',
    description: '',
    sortOrder: Number(newLink.sortOrder || 100) + 10,
    visible: true,
    openInNewTab: true
  })
  toast.show(t('addSuccess'))
}

async function saveGroup(group: Parameters<typeof nav.saveGroup>[0]) {
  await nav.saveGroup(group)
  toast.show(t('saveSuccess'))
}

async function saveLink(link: Parameters<typeof nav.saveLink>[0]) {
  await nav.saveLink(link)
  toast.show(t('saveSuccess'))
}

async function deleteGroup(id: number) {
  await nav.deleteGroup(id)
  toast.show(t('deleteSuccess'))
}

async function deleteLink(id: number) {
  await nav.deleteLink(id)
  toast.show(t('deleteSuccess'))
}
</script>

<template>
  <section class="admin-page">
    <header class="admin-heading">
      <h1>{{ t('links') }}</h1>
      <p>{{ nav.groups.length }} {{ t('groupsCount') }} · {{ flatLinks.length }} {{ t('linksCount') }}</p>
    </header>

    <section class="tool-panel">
      <h2>{{ t('addGroup') }}</h2>
      <div class="form-grid compact">
        <label>{{ t('title') }}<input v-model="newGroup.title" /></label>
        <label>{{ t('sortOrder') }}<input v-model.number="newGroup.sortOrder" type="number" /></label>
        <label class="check-row"><input v-model="newGroup.visible" type="checkbox" />{{ t('visible') }}</label>
        <button class="primary-button" type="button" @click="addGroup"><Plus :size="18" />{{ t('addGroup') }}</button>
      </div>
    </section>

    <section class="tool-panel">
      <h2>{{ t('addLink') }}</h2>
      <div class="form-grid">
        <label>{{ t('group') }}
          <select v-model.number="newLink.groupId">
            <option v-for="group in nav.groups" :key="group.id" :value="group.id">{{ group.title }}</option>
          </select>
        </label>
        <label>{{ t('title') }}<input v-model="newLink.title" /></label>
        <label>{{ t('url') }}<input v-model="newLink.url" placeholder="https://example.com" /></label>
        <label>{{ t('icon') }}<input v-model="newLink.icon" /></label>
        <label>{{ t('description') }}<input v-model="newLink.description" /></label>
        <label>{{ t('sortOrder') }}<input v-model.number="newLink.sortOrder" type="number" /></label>
        <label class="check-row"><input v-model="newLink.visible" type="checkbox" />{{ t('visible') }}</label>
        <label class="check-row"><input v-model="newLink.openInNewTab" type="checkbox" />{{ t('openInNewTab') }}</label>
        <button class="primary-button" type="button" @click="addLink"><Plus :size="18" />{{ t('addLink') }}</button>
      </div>
    </section>

    <section class="tool-panel">
      <h2>{{ t('groupList') }}</h2>
      <div class="admin-list">
        <div v-for="group in nav.groups" :key="group.id" class="admin-row">
          <input v-model="group.title" />
          <input v-model.number="group.sortOrder" type="number" />
          <label class="check-row"><input v-model="group.visible" type="checkbox" />{{ t('visible') }}</label>
          <button class="icon-button" type="button" :title="t('save')" @click="saveGroup(group)"><Save :size="18" /></button>
          <button class="icon-button danger" type="button" :title="t('delete')" @click="deleteGroup(group.id)"><Trash2 :size="18" /></button>
        </div>
      </div>
    </section>

    <section class="tool-panel">
      <h2>{{ t('linkList') }}</h2>
      <div class="admin-list">
        <div v-for="link in flatLinks" :key="link.id" class="admin-row link-row">
          <select v-model.number="link.groupId">
            <option v-for="group in nav.groups" :key="group.id" :value="group.id">{{ group.title }}</option>
          </select>
          <input v-model="link.title" />
          <input v-model="link.url" />
          <input v-model="link.icon" />
          <input v-model="link.description" />
          <input v-model.number="link.sortOrder" type="number" />
          <label class="check-row"><input v-model="link.visible" type="checkbox" />{{ t('visible') }}</label>
          <button class="icon-button" type="button" :title="t('save')" @click="saveLink(link)"><Save :size="18" /></button>
          <button class="icon-button danger" type="button" :title="t('delete')" @click="deleteLink(link.id)"><Trash2 :size="18" /></button>
        </div>
      </div>
    </section>
  </section>
</template>
