<template>
  <div class="page">
    <van-nav-bar title="发布求购" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field v-model="form.book_title" label="需要书籍" placeholder="书名" :rules="[{ required: true, message: '请输入书名' }]" />
        <van-field v-model="form.author" label="作者" placeholder="作者（选填）" />
        <van-field v-model="form.isbn" label="ISBN" placeholder="ISBN（选填）" />
        <van-field v-model="form.expected_price" type="number" label="期望价格" placeholder="期望价格" />
        <van-field label="新旧要求">
          <template #input>
            <van-radio-group v-model="form.condition_requirement" direction="horizontal">
              <van-radio v-for="c in ConditionOptions" :key="c.value" :name="c.value">{{ c.label }}</van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <SubjectPicker v-model="form.subject_category" label="学科分类" placeholder="请选择学科分类" />
        <van-field v-model="form.description" label="备注" type="textarea" rows="3" autosize placeholder="补充说明" maxlength="1000" show-word-limit />
      </van-cell-group>
      <div class="submit-wrap">
        <van-button round block type="primary" native-type="submit" :loading="submitting">发布求购</van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { showSuccessToast } from 'vant'
import { createWish } from '@/api/wish'
import { ConditionOptions } from '@/constants/enums'
import SubjectPicker from '@/components/SubjectPicker.vue'

const router = useRouter()
const submitting = ref(false)
const form = reactive({
  book_title: '',
  author: '',
  isbn: '',
  expected_price: '',
  condition_requirement: 'nine_new',
  subject_category: '',
  description: '',
})

async function onSubmit() {
  if (!form.subject_category) {
    return
  }
  submitting.value = true
  try {
    await createWish({
      book_title: form.book_title,
      author: form.author,
      isbn: form.isbn,
      expected_price: form.expected_price ? Number(form.expected_price) : 0,
      condition_requirement: form.condition_requirement,
      subject_category: form.subject_category,
      description: form.description,
    })
    showSuccessToast('发布成功')
    router.push('/wishes')
  } catch {
    /* toast */
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.submit-wrap {
  padding: 12px 12px 40px;
}
</style>
