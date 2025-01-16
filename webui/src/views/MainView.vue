<template>
  <div class="main-view">
    <div class="conversations-panel">
      <ConversationsView @select-conversation="selectConversation" />
    </div>
    <div class="chat-panel">
      <ChatView v-if="selectedConversationId" :conversationID="selectedConversationId" @conversation-deleted="handleConversationDeleted" />
      <div v-else class="empty-chat">Seleziona una conversazione per iniziare a chattare</div>
    </div>
  </div>
</template>

<script>
import ConversationsView from './ConversationsView.vue';
import ChatView from './ChatView.vue';

export default {
  components: {
    ConversationsView,
    ChatView,
  },
  data() {
    return {
      selectedConversationId: null,
    };
  },
  methods: {
    selectConversation(conversationId) {
      this.selectedConversationId = conversationId;
    },
    handleConversationDeleted(deletedConversationId) {
      // Deseleziona la conversazione se è quella attualmente visualizzata
      if (this.selectedConversationId === deletedConversationId) {
        this.selectedConversationId = null;
      }
    }
  },
};
</script>

<style scoped>
.main-view {
  display: flex;
  height: 100vh; /* Occupa l'intera altezza della finestra */
  background-color: #f0f2f5;
}

.conversations-panel {
  width: 40%; /* Larghezza del pannello delle conversazioni */
  border-right: 1px solid #d1d7db;
  overflow-y: auto; /* Abilita lo scrolling se necessario */
}

.chat-panel {
  width: 60%; /* Larghezza del pannello della chat */
  display: flex;
  flex-direction: column;
}

.empty-chat {
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 1.2em;
  color: #667781;
  height: 100%;
}
</style>