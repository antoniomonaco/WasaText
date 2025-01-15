<template>
    <div class="main-view">
      <div class="conversations-panel">
        <ConversationsView @select-conversation="showChatView" />
      </div>
      <div class="chat-panel">
        <ChatView v-if="selectedConversationId" :conversationID="selectedConversationId" />
        <div v-else class="empty-chat">Seleziona una conversazione</div>
      </div>
    </div>
  </template>
  
  <script>
  import ConversationsView from './ConversationsView.vue';
  import ChatView from './ChatView.vue';
  
  export default {
    components: {
      ConversationsView,
      ChatView
    },
    data() {
      return {
        selectedConversationId: null
      };
    },
    methods: {
      showChatView(conversationId) {
        this.selectedConversationId = conversationId;
      }
    }
  };
  </script>
  
  <style scoped>
  .main-view {
    display: flex;
    height: 100vh; /* Occupa l'intera altezza della finestra */
  }
  
  .conversations-panel {
    width: 30%; /* Larghezza del pannello delle conversazioni */
    border-right: 1px solid #ccc;
    overflow-y: auto; /* Abilita lo scrolling se necessario */
  }
  
  .chat-panel {
    width: 70%; /* Larghezza del pannello della chat */
    display: flex;
    flex-direction: column;
  }
  
  .empty-chat {
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 1.2em;
    color: #999;
    height: 100%;
  }
  </style>