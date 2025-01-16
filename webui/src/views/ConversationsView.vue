<template>
  <div class="conversations-view">
    <div class="search-bar">
      <i class="fas fa-search search-icon"></i>
      <input type="text" placeholder="Cerca" v-model="searchQuery" />
    </div>
    <ul v-if="filteredConversations.length > 0">
      <li
        v-for="conversation in filteredConversations"
        :key="conversation.id"
        class="conversation"
        @click="selectConversation(conversation.id)"
      >
        <img
          v-if="conversation.photoUrl"
          :src="conversation.photoUrl"
          alt="Foto profilo"
          class="conversation-photo"
        />
        <div class="conversation-info">
          <div class="conversation-name-row">
            <span class="conversation-name" v-if="conversation.type === 'group'">
              {{ conversation.name || "Gruppo senza nome" }}
            </span>
            <span class="conversation-name" v-else>
              {{ getOtherParticipant(conversation).username }}
            </span>
            <span class="conversation-time">
              {{ conversation.latestMessage ? formatTimestamp(conversation.latestMessage.timestamp) : "" }}
            </span>
          </div>
          <div class="conversation-preview-row">
            <span class="conversation-preview">
              {{ conversation.latestMessage ? conversation.latestMessage.content : "Nessun messaggio" }}
            </span>
            
            <i v-if="conversation.latestMessage && conversation.latestMessage.sender.id !== currentUserID && conversation.unreadCount > 0" class="fas fa-circle unread-indicator"></i>
            <span v-else-if="conversation.latestMessage && conversation.latestMessage.sender.id === currentUserID" class="fas fa-check sent-indicator"></span>

          </div>
        </div>
      </li>
    </ul>
    <div v-else>Nessuna conversazione trovata.</div>
  </div>
</template>
<script>
export default {
  data() {
    return {
      conversations: [],
      searchQuery: "",
      currentUserID: this.getCurrentUserID(),
    };
  },
  async created() {
    await this.fetchConversations();
  },
  computed: {
    filteredConversations() {
      if (!this.searchQuery) {
        return this.conversations;
      }

      const query = this.searchQuery.toLowerCase();
      return this.conversations.filter((conversation) => {
        const conversationName = (conversation.name || "").toLowerCase();
        const otherParticipantName = (
          this.getOtherParticipant(conversation).username || ""
        ).toLowerCase();
        const latestMessageContent = (
          conversation.latestMessage ? conversation.latestMessage.content : ""
        ).toLowerCase();

        return (
          conversationName.includes(query) ||
          otherParticipantName.includes(query) ||
          latestMessageContent.includes(query)
        );
      });
    },
  },
  methods: {
    getCurrentUserID() {
      const token = localStorage.getItem('authToken');
      return parseInt(token);
    },
    async fetchConversations() {
      console.log("fetch conversations");
      try {
        const token = localStorage.getItem("authToken");
        let authorization = `Bearer ${token}`;
        console.log(authorization);
        const response = await this.$axios.get("/conversations/", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });
        this.conversations = response.data.map(conversation => ({
          ...conversation,
          unreadCount: conversation.messages ? conversation.messages.filter(message => !message.readStatus && message.sender.id !== this.currentUserID).length : 0
        }));
      } catch (error) {
        console.error("Errore nel recupero delle conversazioni:", error);
      }
    },
    getOtherParticipant(conversation) {
      const token = localStorage.getItem("authToken");
      const currentUserID = parseInt(token);
      return conversation.participants.find((p) => p.id !== currentUserID);
    },
    selectConversation(conversationId) {
      this.$emit("select-conversation", conversationId);
    },
    formatTimestamp(timestamp) {
      return new Date(timestamp).toLocaleTimeString([], {
        hour: "2-digit",
        minute: "2-digit",
      });
    },
  },
};
</script>
<style scoped>
@import url("https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.4.2/css/all.min.css");
.conversations-view {
  width: 100%;
  height: 100%;
  background-color: #f8f9fa;
  border-right: 1px solid #e9edef;
  overflow-y: auto;
}
.search-bar {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #f0f2f5;
  border-bottom: 1px solid #d1d7db;
}
.search-icon {
  margin-right: 10px;
  color: #919191;
}
.search-bar input {
  width: 90%;
  border: none;
  background-color: #fff;
  padding: 8px;
  border-radius: 20px;
  outline: none;
  font-size: 14px;
}
.conversation {
  display: flex;
  align-items: center;
  padding: 10px;
  background-color: #fff;
  border-bottom: 1px solid #e9edef;
  cursor: pointer;
  transition: background-color 0.2s;
}
.conversation:hover {
  background-color: #f5f5f5;
}
.conversation-photo {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  margin-right: 15px;
}
.conversation-info {
  flex-grow: 1;
}
.conversation-name-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.conversation-name {
  font-weight: 600;
  color: #111b21;
  font-size: 16px;
}
.conversation-time {
  font-size: 12px;
  color: #667781;
}
.conversation-preview-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.conversation-preview {
  font-size: 14px;
  color: #667781;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 70%;
}
.unread-indicator {
  color: #00a884;
  font-size: 12px;
  margin-left: auto;
}
.sent-indicator {
  color: #667781; /* Colore grigio per il segno di spunta */
  font-size: 14px; /* Dimensione del segno di spunta */
  margin-left: auto; /* Sposta il segno di spunta a destra */
}
</style>