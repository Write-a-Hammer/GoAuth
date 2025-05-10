import { ref, watch } from "vue";
import axios from "axios";

export function useDailyImage(apiBase) {
  const backgroundImage = ref("");
  const isLoading = ref(true);
  const appStyle = ref({
    backgroundImage: "",
    backgroundSize: "cover",
    backgroundPosition: "center",
  });

  const fetchDailyImage = async () => {
    try {
      const response = await axios.get(`${apiBase}/api/bing-image`);
      console.log("API 响应数据:", response.data);
      if (
        response.data &&
        Array.isArray(response.data.images) &&
        response.data.images.length > 0
      ) {
        const imageUrl = `https://www.bing.com${response.data.images[0].url}`;
        const img = new Image();
        img.src = imageUrl;
        img.onload = () => {
          backgroundImage.value = imageUrl;
          isLoading.value = false;
        };
      } else {
        throw new Error("API 返回的数据格式不正确");
      }
    } catch (error) {
      console.error("每日一图加载失败:", error);
      backgroundImage.value = "/fallback-image.jpg";
      isLoading.value = false;
    }
  };

  watch(backgroundImage, (newImage) => {
    appStyle.value.backgroundImage = `url(${newImage})`;
  });

  return {
    backgroundImage,
    isLoading,
    appStyle,
    fetchDailyImage,
  };
}