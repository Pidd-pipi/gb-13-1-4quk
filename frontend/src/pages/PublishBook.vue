<template>
  <div class="page">
    <van-nav-bar :title="isEdit ? '编辑书籍' : '发布闲置书籍'" left-arrow fixed placeholder @click-left="$router.back()" />
    <van-form @submit="onSubmit">
      <van-cell-group inset>
        <van-field v-model="form.title" label="书名" placeholder="请输入书名" :rules="[{ required: true, message: '请输入书名' }]" />
        <van-field v-model="form.author" label="作者" placeholder="请输入作者" />
        <van-field v-model="form.isbn" label="ISBN" placeholder="请输入 ISBN" />
        <van-field v-model="form.course_name" label="课程名" placeholder="如：高等数学" />
        <van-field v-model="form.original_price" type="number" label="原价" placeholder="图书原价" />
        <van-field v-model="form.price" type="number" label="售价" placeholder="出售价格" :rules="[{ required: true, message: '请输入售价' }]" />
        <van-field label="新旧程度" required>
          <template #input>
            <van-radio-group v-model="form.condition" direction="horizontal">
              <van-radio v-for="c in ConditionOptions" :key="c.value" :name="c.value">{{ c.label }}</van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <SubjectPicker v-model="form.subject_category" label="学科分类" placeholder="请选择学科分类" />
        <van-field label="交易方式" required>
          <template #input>
            <van-radio-group v-model="form.trade_type" direction="horizontal">
              <van-radio name="in_person">面交</van-radio>
              <van-radio name="mail">邮寄</van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <van-field v-model="form.campus" label="所在校区" placeholder="如：东校区" />
        <van-field label="支持短借">
          <template #input>
            <van-switch v-model="form.borrowable" @change="onToggleBorrowable" />
            <span class="borrow-hint">教材只用一学期？开启后同学可申请借阅 7 / 14 天</span>
          </template>
        </van-field>
        <van-field v-if="form.borrowable" label="借阅时长" required>
          <template #input>
            <van-radio-group v-model="form.borrow_duration" direction="horizontal">
              <van-radio v-for="d in BorrowDurationOptions" :key="d.value" :name="d.value">{{ d.label }}</van-radio>
            </van-radio-group>
          </template>
        </van-field>
        <van-field v-model="form.description" label="描述" type="textarea" rows="3" autosize placeholder="书籍情况、笔记情况等" maxlength="2000" show-word-limit />
      </van-cell-group>

      <div class="upload-section">
        <div class="upload-title">实物图片（最多 5 张）</div>
        <ImageUploader v-model="images" :max-count="5" />
      </div>

      <div class="submit-wrap">
        <van-button round block type="primary" native-type="submit" :loading="submitting">
          {{ isEdit ? '保存修改' : '发布' }}
        </van-button>
      </div>
    </van-form>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { showSuccessToast, showFailToast } from 'vant'
import { createBook, updateBook, getBook } from '@/api/book'
import { ConditionOptions, BorrowDurationOptions } from '@/constants/enums'
import SubjectPicker from '@/components/SubjectPicker.vue'
import ImageUploader from '@/components/ImageUploader.vue'

const route = useRoute()
const router = useRouter()
const editId = computed(() => (route.query.id ? Number(route.query.id) : 0))
const isEdit = computed(() => editId.value > 0)

const form = reactive({
  title: '',
  author: '',
  isbn: '',
  course_name: '',
  original_price: '',
  price: '',
  condition: 'nine_new',
  subject_category: '',
  trade_type: 'in_person',
  campus: '',
  description: '',
  borrowable: false,
  borrow_duration: 7,
})
const images = ref<string[]>([])
const submitting = ref(false)

function onToggleBorrowable(value: boolean | string | number) {
  if (value && !form.borrow_duration) {
    form.borrow_duration = 7
  }
}

onMounted(async () => {
  if (isEdit.value) {
    try {
      const book = await getBook(editId.value)
      form.title = book.title
      form.author = book.author
      form.isbn = book.isbn
      form.course_name = book.course_name
      form.original_price = String(book.original_price || '')
      form.price = String(book.price)
      form.condition = book.condition
      form.subject_category = book.subject_category
      form.trade_type = book.trade_type
      form.campus = book.campus
      form.description = book.description
      form.borrowable = !!book.borrowable
      form.borrow_duration = book.borrow_duration || 7
      images.value = book.images || []
    } catch {
      /* toast */
    }
  }
})

async function onSubmit() {
  if (!form.subject_category) {
    return
  }
  if (form.borrowable && !form.borrow_duration) {
    showFailToast('请选择借阅时长')
    return
  }
  submitting.value = true
  try {
    const payload = {
      title: form.title,
      author: form.author,
      isbn: form.isbn,
      course_name: form.course_name,
      original_price: form.original_price ? Number(form.original_price) : 0,
      price: Number(form.price),
      condition: form.condition,
      subject_category: form.subject_category,
      trade_type: form.trade_type,
      campus: form.campus,
      description: form.description,
      images: images.value,
      borrowable: form.borrowable,
      borrow_duration: form.borrowable ? form.borrow_duration : 0,
    }
    if (isEdit.value) {
      await updateBook(editId.value, payload)
      showSuccessToast('已保存')
    } else {
      await createBook(payload)
      showSuccessToast('发布成功')
    }
    router.push('/my-books')
  } catch {
    /* toast handled */
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.upload-section {
  background: #fff;
  margin: 12px;
  border-radius: 8px;
  padding: 12px;
}
.upload-title {
  font-size: 14px;
  font-weight: 600;
  margin-bottom: 10px;
}
.borrow-hint {
  margin-left: 10px;
  color: #969799;
  font-size: 12px;
}
.submit-wrap {
  padding: 12px 12px 40px;
}
</style>
