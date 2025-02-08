<!-- ChatView.vue -->
<template>
  <div class="chat-view">
    <!-- Header della chat -->
    <ChatHeader
      :conversationName="conversationName"
      :conversationPhoto="conversationPhoto"
      :is-group="isGroup"
      :participants-count="conversation?.participants?.length || 0"
      @show-participants="showParticipantsModal = true"
      @add-participants="openAddParticipantsModal"
      @edit-group="openEditGroupModal"
      @leave-group="leaveGroup"
    />

    <!-- Area messaggi -->
    <div 
      class="messages-container" 
      ref="messagesContainer"
    >
      <div class="date-separator">{{ getCurrentDate() }}</div>
      
      <!-- Lista messaggi -->
      <div class="messages-list">
        <MessageItem
          v-for="message in messages"
          :key="message.id"
          :message="message"
          :current-user-id="currentUserId"
          :is-group="isGroup"
          :is-message-menu-open="selectedMessage === message.id"
          @open-message-menu="openMessageMenu"
          @comment-message="commentMessage"
          @reply-message="replyToMessage"
          @show-comments="showComments"
          @forward-message="forwardMessage"
          @delete-message="deleteMessage"
        />
      </div>
    </div>

    <!-- Area input messaggi -->
    <ChatInput
      :replying-to="replyingTo"
      :replying-message="replyingToMessage"
      @send-message="sendMessage"
      @cancel-reply="cancelReply"
    />

    <!-- Modal Commenti -->
    <BaseModal 
      v-if="showCommentsModal"
      title="Commenti"
      @close="closeCommentsModal"
    >
      <CommentsModal
        :comments="currentComments"
        :loading="isLoadingComments"
        :error="commentsError"
        :current-user-id="currentUserId"
        @submit="submitComment"
        @delete-comment="deleteComment"
      >
        <template #comment-input>
          <input 
            v-model="newComment" 
            type="text"
            placeholder="Scrivi un commento..."
            @keyup.enter="submitComment"
            :disabled="isLoadingComments"
          />
        </template>
      </CommentsModal>
    </BaseModal>

    <!-- Modal Emoji Picker per Commento -->
    <div 
      v-if="showEmojiPickerModal" 
      class="emoji-picker-overlay"
      @click.self="closeEmojiPickerModal"
    >
      <div class="emoji-picker-container">
        <div class="emoji-picker-header">
          <h3>Aggiungi Reazione</h3>
          <button 
            class="close-button" 
            @click="closeEmojiPickerModal"
          >
            <i class="fas fa-times"></i>
          </button>
        </div>
        <EmojiPicker 
          @select="submitEmojiComment"
        />
      </div>
    </div>

    <!-- Resto dei modal precedenti -->
    <BaseModal 
      v-if="showForwardModal"
      title="Inoltra messaggio"
      @close="closeForwardModal"
    >
      <ForwardModal
        :conversations="forwardingConversations"
        :current-user-id="currentUserId"
        :loading="isLoadingConversations"
        :exclude-conversation-id="conversationID"
        @select="confirmForward"
      />
    </BaseModal>

    <BaseModal 
      v-if="showParticipantsModal"
      title="Partecipanti del gruppo"
      @close="showParticipantsModal = false"
    >
      <ParticipantsModal
        :participants="conversation?.participants || []"
        :current-user-id="currentUserId"
        :is-group="isGroup"
        :loading="isLoadingParticipants"
        @close="showParticipantsModal = false" 
        @remove-participant="removeParticipant"
        @add-participants="openAddParticipantsModal"
      />
    </BaseModal>

    <BaseModal 
      v-if="showEditGroupModal"
      title="Modifica gruppo"
      @close="showEditGroupModal = false"
    >
      <GroupEditModal
        :initial-name="conversation?.name || ''"
        :initial-photo-url="conversation?.photoUrl || ''"
        @save="saveGroupChanges"
      />
    </BaseModal>

    <BaseModal 
      v-if="showAddParticipantsModal"
      title="Aggiungi partecipanti"
      @close="closeAddParticipantsModal"
    >
      <div class="add-participants-content">
        <input 
          type="text"
          v-model="userSearchQuery"
          placeholder="Cerca utenti..."
          class="search-input"
        />
        <div class="users-list">
          <div 
            v-for="user in filteredUsers" 
            :key="user.id"
            class="user-item"
            @click="addParticipant(user)"
          >
            <img 
              :src="user.photoUrl || '/default-avatar.jpeg'" 
              :alt="user.username"
              class="user-avatar"
            />
            <span class="user-name">{{ user.username }}</span>
          </div>
        </div>
        <div v-if="filteredUsers.length === 0" class="no-participants">
          Non ci sono altri utenti da aggiungere.
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script>
import ChatHeader from '@/components/ChatHeader.vue'
import MessageItem from '@/components/MessageItem.vue'
import ChatInput from '@/components/ChatInput.vue'
import BaseModal from '@/components/BaseModal.vue'
import CommentsModal from '@/components/CommentsModal.vue'
import ForwardModal from '@/components/ForwardModal.vue'
import ParticipantsModal from '@/components/ParticipantsModal.vue'
import GroupEditModal from '@/components/GroupEditModal.vue'
import EmojiPicker from '@/components/EmojiPicker.vue'

export default {
  name: 'ChatView',
  components: {
    ChatHeader,
    MessageItem,
    ChatInput,
    BaseModal,
    CommentsModal,
    ForwardModal,
    ParticipantsModal,
    GroupEditModal,
    EmojiPicker
  },
  props: {
    conversationID: {
      type: Number,
      required: true
    }
  },
  data() {
    return {
      // Stato della conversazione
      conversation: null,
      messages: [],
      currentUserId: parseInt(localStorage.getItem('authToken')),

      // Stato dei messaggi
      selectedMessage: null,
      replyingTo: null,

      // Modal e loro stati
      showCommentsModal: false,
      showForwardModal: false,
      showParticipantsModal: false,
      showEditGroupModal: false,
      showAddParticipantsModal: false,
      showEmojiPickerModal: false,

      // Stato dei commenti
      currentComments: [],
      isLoadingComments: false,
      commentsError: null,
      selectedMessageForComments: null,

      // Stato dell'inoltro
      forwardingConversations: [],
      isLoadingConversations: false,
      messageToForward: null,

      // Stato dei partecipanti
      isLoadingParticipants: false,
      userSearchQuery: '',
      users: [],

      pollingInterval : null,
      isNearBottom: true,
      scrollThreshold: 100, // pixel dal basso che permettono la chiamata a scrollToBottom()

      replyingToMessage: null,
      newComment: ''
    }
  },
  computed: {
    conversationName() {
      if (!this.conversation) return 'Caricamento...'
      
      if (this.isGroup) {
        return this.conversation.name || 'Gruppo senza nome'
      }
      
      return this.getOtherParticipant()?.username || 'Chat'
    },
    conversationPhoto() {
      if (!this.conversation) return '/default-avatar.jpeg'
      
      if (this.isGroup) {
        return this.conversation.photoUrl || '/default-avatar.jpeg'
      }
      
      return this.getOtherParticipant()?.photoUrl || '/default-avatar.jpeg'
    },
    isGroup() {
      return this.conversation?.type === 'group'
    },
    filteredUsers() {
      if (!this.userSearchQuery) {
        return this.users.filter(user => 
          !this.conversation.participants.some(participant => participant.id === user.id)
        );
      }

      const query = this.userSearchQuery.toLowerCase()
      return this.users.filter(user => 
        user.username.toLowerCase().includes(query) &&
        !this.conversation.participants.some(p => p.id === user.id)
      )
    }
  },
  methods: {
    // Metodi di utilità
    getOtherParticipant() {
      return this.conversation?.participants?.find(p => p.id !== this.currentUserId)
    },
    getCurrentDate() {
      return new Date().toLocaleDateString('it-IT', {
        weekday: 'long', 
        day: 'numeric', 
        month: 'long', 
        year: 'numeric'
      })
    },

    // Metodi di gestione conversazione
    async fetchConversation() {
      try {
        const response = await this.$axios.get(
          `/conversations/${this.conversationID}`, 
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        this.conversation = response.data.conversation;
        const newMessages = response.data.messages || [];
        
        // Carica i commenti per ogni messaggio
        const messagesWithComments = await Promise.all(
          newMessages.map(async (message) => {
            try {
              const commentsResponse = await this.$axios.get(
                `/conversations/${this.conversationID}/messages/${message.id}/comments`,
                {
                  headers: {
                    Authorization: `Bearer ${localStorage.getItem('authToken')}`
                  }
                }
              );
              return {
                ...message,
                comments: commentsResponse.data
              };
            } catch (error) {
              console.error(`Error fetching comments for message ${message.id}:`, error);
              return {
                ...message,
                comments: []
              };
            }
          })
        );

        const hasNewMessages = !this.messages || 
          messagesWithComments.length !== this.messages.length ||
          JSON.stringify(messagesWithComments) !== JSON.stringify(this.messages);

        if (hasNewMessages) {
          this.messages = messagesWithComments;
          this.$nextTick(() => {
            this.scrollToBottom();
            this.markMessagesAsRead();
          });
        }
      } catch (error) {
        console.error('Errore nel recupero della conversazione:', error);
      }
    },
    async updateMessageComments(messageId) {
      try {
        const commentsResponse = await this.$axios.get(
          `/conversations/${this.conversationID}/messages/${messageId}/comments`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        this.messages = this.messages.map(message => {
          if (message.id === messageId) {
            return {
              ...message,
              comments: commentsResponse.data
            };
          }
          return message;
        });
      } catch (error) {
        console.error(`Error updating comments for message ${messageId}:`, error);
      }
    },

    async replyToMessage(message) {
      this.replyingTo = {
        id: message.id
      };
      try {
        const response = await this.$axios.get(
          `/conversations/${this.conversationID}/messages/${message.id}`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        this.replyingToMessage = response.data;
      } catch (error) {
        console.error('Error fetching reply message:', error);
        this.replyingTo = null;
        this.replyingToMessage = null;
      }
      this.selectedMessage = null;
    },

    // Metodi di gestione messaggi
    async sendMessage(messagePayload) {
      try {
        const response = await this.$axios.post(
          `/conversations/${this.conversationID}`,
          {
            type: messagePayload.type || 'text',
            content: messagePayload.type === 'media' ? messagePayload.content : messagePayload.text,
            replyTo: this.replyingToMessage ? this.replyingToMessage.id : 0
          },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );

        if (response.data) {
          const newMessage = {
            ...response.data,
            replyTo: this.replyingToMessage ? this.replyingToMessage.id : 0
          };
          
          this.messages = [...this.messages, newMessage];
          this.$nextTick(() => {
            this.scrollToBottom();
          });
        }

        this.replyingTo = null;
        this.replyingToMessage = null;
      } catch (error) {
        console.error('Error sending message:', error);
      }
    },

    async markMessagesAsRead() {
      try {
        await this.$axios.put(
          `/conversations/${this.conversationID}`,
          {},
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )
      } catch (error) {
        console.error('Errore nel marcare i messaggi come letti:', error)
      }
    },

    // Metodi di gestione menu contestuale
    openMessageMenu(message, event) {
      // Chiude il menu se è già aperto per lo stesso messaggio
      if (this.selectedMessage === message.id) {
        this.selectedMessage = null
      } else {
        this.selectedMessage = message.id
      }
      
      // Chiude altri modal aperti
      this.showCommentsModal = false
      this.showForwardModal = false
      
      // Per non chiudere il menu immediatamente
      event.stopPropagation()
    },
    mounted() {
      document.addEventListener('click', this.handleOutsideClick)
      const container = this.$refs.messagesContainer;
      if (container) {
        container.addEventListener('scroll', this.handleScroll);
      }
    },

    beforeUnmount() {
      document.removeEventListener('click', this.handleOutsideClick)
      if (container) {
        container.removeEventListener('scroll', this.handleScroll);
      }
    },

    handleOutsideClick(event) {
      // Se il menu è aperto e si clicca fuori, chiudilo
      if (this.selectedMessage && !event.target.closest('.message-wrapper')) {
        this.selectedMessage = null
      }
    },

    // Metodi di gestione commenti
    async showComments(message) {
      this.selectedMessageForComments = message
      this.showCommentsModal = true
      this.isLoadingComments = true
      this.commentsError = null

      try {
        const response = await this.$axios.get(
          `/conversations/${this.conversationID}/messages/${message.id}/comments`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )
        this.currentComments = response.data
      } catch (error) {
        console.error('Errore nel recupero dei commenti:', error)
        this.commentsError = 'Impossibile caricare i commenti'
      } finally {
        this.isLoadingComments = false
      }
    },
    async submitComment(commentText) {
      if (!commentText.trim() || !this.selectedMessageForComments) return

      try {
        const response = await this.$axios.post(
          `/conversations/${this.conversationID}/messages/${this.selectedMessageForComments.id}/comments`,
          { content: commentText.trim() },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )

        if (response.data) {
          this.currentComments = [...this.currentComments, response.data]
          this.newComment = ''
        }
      } catch (error) {
        console.error('Errore nell\'invio del commento:', error)
      }
    },

    async deleteComment(comment) {
      try {
        await this.$axios.delete(
          `/conversations/${this.conversationID}/messages/${this.selectedMessageForComments.id}/comments/${comment.id}`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )
        
        await this.showComments(this.selectedMessageForComments)
      } catch (error) {
        console.error('Errore nell\'eliminazione del commento:', error)
      }
    },

    closeCommentsModal() {
      this.showCommentsModal = false
      this.selectedMessageForComments = null
      this.currentComments = []
    },

        // Metodo per eliminare un messaggio
    async deleteMessage(message) {
      try {
        await this.$axios.delete(
          `/conversations/${this.conversationID}/messages/${message.id}`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )
        
        // Ricarica la conversazione dopo l'eliminazione
        await this.fetchConversation()
      } catch (error) {
        console.error('Errore nell\'eliminazione del messaggio:', error)
      }
    },

    // Metodo per commentare un messaggio
    commentMessage(message) {
      this.messageForEmojiComment = message
      this.showEmojiPickerModal = true
    },

    closeEmojiPickerModal() {
      this.showEmojiPickerModal = false
      this.messageForEmojiComment = null
    },

    async submitEmojiComment(emoji) {
      if (!this.messageForEmojiComment) return

      try {
        const response = await this.$axios.post(
          `/conversations/${this.conversationID}/messages/${this.messageForEmojiComment.id}/comments`,
          { content: emoji },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )

        this.closeEmojiPickerModal()
      } catch (error) {
        console.error('Errore nell\'invio del commento:', error)
      }
    },

    // Metodi di inoltro
    async forwardMessage(message) {
      this.messageToForward = message
      this.showForwardModal = true
      this.isLoadingConversations = true

      try {
        const [conversationsRes, usersRes] = await Promise.all([
          this.$axios.get('/conversations/', {
            headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` }
          }),
          this.$axios.get('/users/', {
            headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` }
          })
        ]);

        // Filtro gli utenti che sono già nelle conversazioni
        const existingUserIds = new Set(
          conversationsRes.data
            .filter(conv => conv.type === 'chat')
            .map(conv => conv.participants.find(p => p.id !== this.currentUserId)?.id)
        );

        const newUsers = usersRes.data.filter(
          user => user.id !== this.currentUserId && !existingUserIds.has(user.id)
        );

        this.forwardingConversations = [
          ...conversationsRes.data,
          ...newUsers.map(user => ({
            isNewChat: true,
            type: 'chat',
            participants: [user],
            name: user.username,
            photoUrl: user.photoUrl
          }))
        ];
      } catch (error) {
        console.error('Errore nel recupero delle conversazioni:', error)
      } finally {
        this.isLoadingConversations = false
      }
    },

    async confirmForward(target) {
      try {
        let targetConversationId = target;
        
        // Se il target è un utente, creo prima la conversazione
        if (typeof target === 'object' && target.isNewChat) {
          const createConvRes = await this.$axios.post(
            '/conversations/',
            {
              type: 'chat',
              participants: [target.participants[0].id]
            },
            {
              headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` }
            }
          );
          targetConversationId = createConvRes.data.conversation_id;
        }

        // Faccio il forward 
        await this.$axios.post(
          `/conversations/${this.conversationID}/messages/${this.messageToForward.id}`,
          { id: targetConversationId },
          {
            headers: { Authorization: `Bearer ${localStorage.getItem('authToken')}` }
          }
        )

        this.closeForwardModal()
      } catch (error) {
        console.error('Errore nell\'inoltrare il messaggio:', error)
      }
    },

    closeForwardModal() {
      this.showForwardModal = false
      this.messageToForward = null
      this.forwardingConversations = []
    },

    // Metodi di gestione gruppo
    cancelReply() {
      this.replyingTo = null
    },

    async leaveGroup() {
      if (!confirm('Sei sicuro di voler uscire dal gruppo?')) return

      try {
        await this.$axios.delete(
          `/conversations/${this.conversationID}/participants`,
          {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          },
          data: { id: this.currentUserId}
        }
      )
        
        this.$emit('conversation-deleted', this.conversationID)
        window.location.href = '/'
      } catch (error) {
        console.error('Errore durante l\'uscita dal gruppo:', error)
      }
    },

    // Metodi di gestione partecipanti
    async fetchUsers() {
      try {
        const response = await this.$axios.get('/users/', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        })
        this.users = response.data
      } catch (error) {
        console.error('Errore nel recupero degli utenti:', error)
      }
    },

    openAddParticipantsModal() {
      this.showAddParticipantsModal = true
      this.userSearchQuery = ''
      this.fetchUsers()
    },

    closeAddParticipantsModal() {
      this.showAddParticipantsModal = false
      this.userSearchQuery = ''
    },

    async addParticipant(user) {
      try {
        await this.$axios.post(
          `/conversations/${this.conversationID}/participants`,
          { id: user.id },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )
        
        await this.fetchConversation()
        this.showAddParticipantsModal = false
      } catch (error) {
        console.error('Errore nell\'aggiunta del partecipante:', error)
      }
    },

    async removeParticipant(participant) {
      try {
        console.log("participant id: ",participant.id)
        console.log("conversation id: ",this.conversationID)
        await this.$axios.delete(
          `/conversations/${this.conversationID}/participants`,
        {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          },
          data: { id: participant.id }
        }
      )
        await this.fetchConversation()
      } catch (error) {
        console.error('Errore nella rimozione del partecipante:', error)
      }
    },

    // Metodi di modifica gruppo
    openEditGroupModal() {
      this.showEditGroupModal = true
    },

    async saveGroupChanges(groupData) {
      try {
        // Salva il nome del gruppo
        await this.$axios.put(
          `/conversations/${this.conversationID}/name`,
          { name: groupData.name },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )

        // Salva la foto del gruppo (se fornita)
        if (groupData.photoUrl) {
          await this.$axios.put(
            `/conversations/${this.conversationID}/photo`,
            { photoUrl: groupData.photoUrl },
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('authToken')}`
              }
            }
          )
        }

        await this.fetchConversation()
        this.showEditGroupModal = false
      } catch (error) {
        console.error('Errore nel salvare le modifiche del gruppo:', error)
      }
    },

    // Metodi di navigazione e scrolling
    handleScroll() {
      const container = this.$refs.messagesContainer;
      if (!container) return;
      
      // Calcola se simao vicini al fondo della chat
      const distanceFromBottom = 
        container.scrollHeight - 
        (container.scrollTop + container.clientHeight);
      
      this.isNearBottom = distanceFromBottom <= this.scrollThreshold;
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer;
        if (container && this.isNearBottom) {
          container.scrollTop = container.scrollHeight;
        }
      });
    },
    startPolling() {
        this.pollingInterval = setInterval(this.fetchConversation, 1000);
      },

    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
      }
    }
  },
  
  async created() {
    this.messages = []; // Inizializzo i messaggi per evitare erroriu nelle conversazioni vuote
    await this.fetchConversation();
    this.startPolling();
  },

  beforeUnmount() {
    this.stopPolling();
  },
  
  // Watcher per gestire cambiamenti nella conversazione
  watch: {
    conversationID: {
      immediate: true,
      async handler(newId) {
        if (newId) {
          await this.fetchConversation();
        }
      }
    }
  }
}
</script>

<style scoped>
.chat-view {
  display: flex;
  flex-direction: column;
  height: 100vh;
  width: 100%;
  background-color: #0b141a;
  position: relative;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px 5% 20px;
  background-color: #0b141a;
  background-image: url("data:image/svg+xml,%3Csvg width='64' height='64' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill='%23192127' d='M8 16h48v32H8z' fill-opacity='.05'/%3E%3C/svg%3E");
}

.date-separator {
  text-align: center;
  color: #8696a0;
  font-size: 12.5px;
  padding: 8px 0;
  position: sticky;
  top: 0;
  background-color: rgba(11, 20, 26, 0.9);
  z-index: 1;
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.add-participants-content {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.search-input {
  width: 100%;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  color: #e9edef;
  padding: 12px;
  font-size: 15px;
  outline: none;
}

.search-input::placeholder {
  color: #8696a0;
}

.users-list {
  max-height: 300px;
  overflow-y: auto;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-radius: 8px;
}

.user-item:hover {
  background-color: #2a3942;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
}

.user-name {
  color: #e9edef;
  font-size: 16px;
}

.no-participants {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    flex: 1;
    color: #8696a0;
    text-align: center;
    padding: 20px;
  }

.emoji-picker-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.emoji-picker-container {
  background-color: #1f2c34;
  border-radius: 12px;
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  max-height: 80vh;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.4);
  overflow: hidden;
}

.emoji-picker-header {
  padding: 16px 20px;
  border-bottom: 1px solid #2a3942;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #202c33;
}

.emoji-picker-header h3 {
  color: #e9edef;
  margin: 0;
  font-size: 16px;
}

.close-button {
  background: none;
  border: none;
  color: #aebac1;
  font-size: 24px;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
  border-radius: 50%;
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.close-button:hover {
  background-color: #384147;
}
</style>