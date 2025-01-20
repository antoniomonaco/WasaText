<!-- EmojiPicker.vue -->
<template>
    <div class="emoji-picker">
      <div class="emoji-categories">
        <button 
          v-for="category in categories" 
          :key="category.name"
          class="category-button"
          :class="{ active: currentCategory === category.name }"
          @click="currentCategory = category.name"
        >
          <i :class="category.icon"></i>
        </button>
      </div>
      
      <div class="emoji-grid">
        <button 
          v-for="emoji in currentEmojis" 
          :key="emoji"
          class="emoji-button"
          @click="$emit('select', emoji)"
        >
          {{ emoji }}
        </button>
      </div>
    </div>
</template>
  
<script>
    export default {
    name: 'EmojiPicker',
    data() {
        return {
        currentCategory: 'smileys',
        categories: [
            { name: 'smileys', icon: 'fas fa-smile', emojis: ['😀', '😁', '😂', '🤣', '😃', '😄', '😅', '😆', '😉', '😊', '😋', '😎', '😍', '😘', '😗', '😙', '😚', '🙂', '🤗', '🤩'] },
            { name: 'gestures', icon: 'fas fa-hand-peace', emojis: ['👍', '👎', '👌', '✌️', '🤞', '👊', '✋', '🤚', '🖐️', '🖖', '👋', '🤙', '💪', '🖕', '✍️', '🙏', '🤝', '💅', '👂', '👃'] },
            { name: 'hearts', icon: 'fas fa-heart', emojis: ['❤️', '🧡', '💛', '💚', '💙', '💜', '🖤', '💔', '❣️', '💕', '💞', '💓', '💗', '💖', '💘', '💝', '💟', '☮️', '✝️', '☪️'] },
            { name: 'nature', icon: 'fas fa-leaf', emojis: ['🌺', '🌹', '🌷', '🌸', '🌼', '🌻', '🌞', '🌝', '🌛', '🌜', '🌚', '🌕', '🌖', '🌗', '🌘', '🌑', '🌒', '🌓', '🌔', '🌙'] },
            { name: 'objects', icon: 'fas fa-lightbulb', emojis: ['⭐️', '🌟', '✨', '⚡️', '🔥', '💫', '☄️', '💥', '🌈', '☀️', '🌤️', '⛅️', '🌥️', '☁️', '🌦️', '🌧️', '⛈️', '🌩️', '🌨️', '💨'] }
        ]
        }
    },
    computed: {
        currentEmojis() {
        const category = this.categories.find(c => c.name === this.currentCategory)
        return category ? category.emojis : []
        }
    }
    }
</script>
  
<style scoped>
.emoji-picker {
background-color: #2a3942;
border-radius: 8px;
overflow: hidden;
width: 100%;
max-width: 360px;
}

.emoji-categories {
display: flex;
background-color: #202c33;
padding: 8px;
border-bottom: 1px solid #374045;
}

.category-button {
background: none;
border: none;
color: #8696a0;
padding: 8px;
cursor: pointer;
border-radius: 4px;
transition: all 0.2s;
}

.category-button:hover,
.category-button.active {
background-color: #374045;
color: #e9edef;
}

.emoji-grid {
display: grid;
grid-template-columns: repeat(8, 1fr);
gap: 4px;
padding: 8px;
max-height: 200px;
overflow-y: auto;
}

.emoji-button {
background: none;
border: none;
font-size: 24px;
padding: 8px;
cursor: pointer;
border-radius: 4px;
transition: background-color 0.2s;
}

.emoji-button:hover {
background-color: #374045;
}

/* Scrollbar styles */
.emoji-grid::-webkit-scrollbar {
width: 6px;
}

.emoji-grid::-webkit-scrollbar-track {
background: transparent;
}

.emoji-grid::-webkit-scrollbar-thumb {
background: #374045;
border-radius: 3px;
}

.emoji-grid::-webkit-scrollbar-thumb:hover {
background: #455A64;
}
</style>