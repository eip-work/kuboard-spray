import { useI18n } from 'vue-i18n';
import { markRaw }from "vue";

export default {
  data () {
    const { t, locale } = useI18n({
      inheritLocale: true,
      useScope: 'local'
    })
    let temp = markRaw(t);
    return { t: temp, locale }
  }

}