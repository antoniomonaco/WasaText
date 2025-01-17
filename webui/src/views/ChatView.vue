<template>
  <div class="chat-view">
    <div class="chat-header">
      <div class="chat-header-info">
        <img 
          :src="conversationPhoto" 
          :alt="conversationName"
          class="chat-avatar"
          @click="showParticipantsModal = true"
        />
        <div class="chat-details">
          <div class="chat-name">{{ conversationName }}</div>
          <div class="chat-status">
            {{ isGroup ? `${conversation?.participants?.length || 0} partecipanti` : 'online' }}
          </div>
        </div>
      </div>
      
      <!-- Pulsanti per le operazioni di gruppo -->
      <div class="chat-header-actions" v-if="isGroup">
        <button class="icon-button" @click="toggleGroupMenu">
          <i class="fas fa-ellipsis-v"></i>
        </button>
        
        <!-- Menu contestuale del gruppo -->
        <div v-if="showGroupMenu" class="group-menu">
          <div class="menu-item" @click="openAddParticipantsModal">
            <i class="fas fa-user-plus"></i>
            Aggiungi partecipanti
          </div>
          <div class="menu-item" @click="showParticipantsModal = true">
            <i class="fas fa-users"></i>
            Visualizza partecipanti
          </div>
          <div class="menu-item" @click="openEditGroupModal">
            <i class="fas fa-edit"></i>
            Modifica gruppo
          </div>
          <div class="menu-item delete" @click="leaveGroup">
            <i class="fas fa-sign-out-alt"></i>
            Esci dal gruppo
          </div>
        </div>
      </div>
    </div>

    <!-- Area messaggi -->
    <div class="messages-container" ref="messagesContainer">
      <div class="date-separator">{{ getCurrentDate() }}</div>
      
      <!-- Lista messaggi -->
      <div class="messages-list">
        <div 
          v-for="message in messages" 
          :key="message.id"
          class="message-wrapper"
          :class="{
            'sent': message.sender.id === currentUserId,
            'received': message.sender.id !== currentUserId
          }"
        >
          <div 
            class="message"
            @contextmenu.prevent="openMessageMenu(message)"
          >
            <!-- Nome mittente (solo per gruppi) -->
            <div v-if="isGroup && message.sender.id !== currentUserId" class="message-sender">
              {{ message.sender.username }}
            </div>

            <!-- Contenuto messaggio -->
            <div class="message-content">
              <!-- Messaggio citato -->
              <div v-if="message.replyTo" class="quoted-message">
                <div class="quoted-sender">
                  {{ getMessageSender(message.replyTo) }}
                </div>
                <div class="quoted-content">
                  {{ getMessageContent(message.replyTo) }}
                </div>
              </div>

              <!-- Testo messaggio -->
              <div v-if="message.type === 'text'" class="message-text">
                {{ message.content }}
              </div>

              <!-- Media -->
              <img 
                v-else-if="message.type === 'media'"
                :src="message.content"
                alt="Media"
                class="message-media"
              />

              <!-- Badge commenti -->
              <div v-if="message.comments?.length" class="comments-badge">
                <i class="fas fa-comment"></i>
                {{ message.comments.length }}
              </div>

              <!-- Metadata (ora e stato) -->
              <div class="message-meta">
                <span class="message-time">
                  {{ formatTime(message.timestamp) }}
                </span>
                <span v-if="message.sender.id === currentUserId" class="message-status">
                  <i :class="[
                    message.status === 'read' ? 'fas fa-check-double read' : 'fas fa-check'
                  ]"></i>
                </span>
              </div>
            </div>
          </div>

          <!-- Menu contestuale del messaggio -->
          <div 
            v-if="selectedMessage === message.id" 
            class="message-menu"
            :class="{'menu-sent': message.sender.id === currentUserId}"
          >
            <div class="menu-item" @click="commentMessage(message)">
              <i class="fas fa-comment"></i>
              Commenta
            </div>
            <div class="menu-item" @click="showComments(message)">
              <i class="fas fa-comments"></i>
              Visualizza commenti
            </div>
            <div class="menu-item" @click="forwardMessage(message)">
              <i class="fas fa-share"></i>
              Inoltra
            </div>
            <div 
              v-if="message.sender.id === currentUserId" 
              class="menu-item delete"
              @click="deleteMessage(message)"
            >
              <i class="fas fa-trash"></i>
              Elimina
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Input area -->
    <div class="input-area">
      <!-- Area risposta -->
      <div v-if="replyingTo" class="reply-container">
        <div class="reply-content">
          <div class="reply-sender">
            Risposta a {{ replyingTo.sender.username }}
          </div>
          <div class="reply-text">
            {{ getMessageContent(replyingTo) }}
          </div>
        </div>
        <button class="icon-button" @click="cancelReply">
          <i class="fas fa-times"></i>
        </button>
      </div>

      <!-- Input messaggi -->
      <div class="input-wrapper">
        <input 
          v-model="newMessage" 
          type="text"
          placeholder="Scrivi un messaggio"
          @keyup.enter="sendMessage"
          class="message-input"
        />
        <button 
          class="send-button"
          :disabled="!newMessage.trim()"
          @click="sendMessage"
        >
          <i class="fas fa-paper-plane"></i>
        </button>
      </div>
    </div>

    <!-- Modal per i commenti -->
    <div v-if="showCommentsModal" class="modal">
      <div class="modal-content comments-modal">
        <div class="modal-header">
          <h3>Commenti</h3>
          <button class="close-button" @click="closeCommentsModal">×</button>
        </div>
        <div class="modal-body">
          <!-- Loading state -->
          <div v-if="isLoadingComments" class="comments-loading">
            <div class="loading-spinner"></div>
            <span>Caricamento commenti...</span>
          </div>

          <!-- Error state -->
          <div v-else-if="commentsError" class="comments-error">
            {{ commentsError }}
          </div>

          <!-- Comments list -->
          <div v-else class="comments-list">
            <div v-if="!currentComments?.length" class="no-comments">
              Nessun commento. Sii il primo a commentare!
            </div>
            <div v-else v-for="comment in currentComments" :key="comment.id" class="comment">
              <div class="comment-header">
                <span class="comment-sender">{{ comment.sender.username }}</span>
                <span class="comment-time">{{ formatTime(comment.timestamp) }}</span>
              </div>
              <div class="comment-content">{{ comment.content }}</div>
              <button 
                v-if="comment.sender.id === currentUserId"
                class="delete-comment"
                @click="deleteComment(comment)"
              >
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
          
          <!-- Input nuovo commento -->
          <div class="new-comment">
            <input 
              v-model="newComment" 
              type="text"
              placeholder="Scrivi un commento..."
              @keyup.enter="submitComment"
              :disabled="isLoadingComments"
            />
            <button 
              @click="submitComment"
              :disabled="!newComment.trim() || isLoadingComments"
            >
              <i class="fas fa-paper-plane"></i>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal per l'inoltro -->
    <div v-if="showForwardModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Inoltra messaggio</h3>
          <button class="close-button" @click="closeForwardModal">×</button>
        </div>
        <div class="modal-body">
          <div v-if="forwardingConversations.length === 0" class="no-conversations">
            Nessuna conversazione disponibile per l'inoltro.
          </div>
          <div class="conversations-list">
            <div 
              v-for="conv in forwardingConversations" 
              :key="conv.id"
              class="conversation-item"
              @click="confirmForward(conv.id)"
            >
              <img 
                :src="conv.photoUrl || '/api/placeholder/40/40'" 
                :alt="getConversationName(conv)"
                class="conversation-avatar"
              />
              <div class="conversation-info">
                <span class="conversation-name">{{ getConversationName(conv) }}</span>
                <span class="conversation-type">
                  {{ conv.type === 'group' ? 'Gruppo' : 'Chat privata' }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- Modal per aggiungere partecipanti -->
    <div v-if="showAddParticipantsModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Aggiungi partecipanti</h3>
          <button class="close-button" @click="closeAddParticipantsModal">×</button>
        </div>
        <div class="modal-body">
          <input 
            type="text"
            class="search-input"
            v-model="userSearchQuery"
            placeholder="Cerca utenti..."
          />
          <div class="users-list">
            <div 
              v-for="user in filteredUsers" 
              :key="user.id"
              class="user-item"
              @click="addParticipant(user)"
            >
              <img 
                :src="user.photoUrl || '/api/placeholder/40/40'" 
                alt="User avatar" 
                class="user-avatar"
              />
              <span class="user-name">{{ user.username }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal per modificare il gruppo -->
    <div v-if="showEditGroupModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3>Modifica gruppo</h3>
          <button class="close-button" @click="closeEditGroupModal">×</button>
        </div>
        <div class="modal-body">
          <input 
            type="text"
            class="modal-input"
            v-model="editGroupName"
            placeholder="Inserisci il nome del gruppo"
          />
          <button 
            class="modal-button"
            :disabled="!editGroupName.trim()"
            @click="saveGroupChanges"
          >
            Salva modifiche
          </button>
        </div>
      </div>
    </div>

    <!-- Additional modal for participants -->
    <div v-if="showParticipantsModal" class="modal">
      <div class="modal-content participants-modal">
        <div class="modal-header">
          <h3>Partecipanti del gruppo</h3>
          <button class="close-button" @click="showParticipantsModal = false">×</button>
        </div>
        <div class="modal-body">
          <div class="participants-list">
            <div 
              v-for="participant in conversation?.participants" 
              :key="participant.id" 
              class="participant-item"
            >
              <img 
                :src="participant.photoUrl || '/api/placeholder/40/40'" 
                :alt="participant.username"
                class="participant-avatar"
              />
              <div class="participant-info">
                <span class="participant-name">{{ participant.username }}</span>
                <span 
                  v-if="participant.id === currentUserId" 
                  class="participant-you"
                >
                  (Tu)
                </span>
              </div>
              <button 
                v-if="isGroup && participant.id !== currentUserId" 
                class="remove-participant-btn"
                @click="removeParticipant(participant)"
              >
                <i class="fas fa-times"></i>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    conversationID: {
      type: Number,
      required: true
    }
  },

  data() {
    return {
      conversation: null,
      messages: [],
      newMessage: '',
      currentUserId: parseInt(localStorage.getItem('authToken')),
      selectedMessage: null,
      replyingTo: null,
      showForwardModal: false,
      messageToForward: null,
      forwardingConversations: [],
      pollingInterval: null,
      showCommentsModal: false,
      currentComments: [],
      newComment: '',
      selectedMessageForComments: null,
      isLoadingComments: false,
      commentsError: null,
      showGroupMenu: false,
      showParticipantsModal: false,
      showAddParticipantsModal: false,
      showEditGroupModal: false,
      userSearchQuery: '',
      users: [],
      editGroupName: '',
      editGroupPhoto: '',
      isLoadingUsers: false,
    }
  },

  computed: {
    conversationName() {
      if (this.conversation?.type === 'group') {
        return this.conversation.name || 'Gruppo senza nome'
      }
      return this.getOtherParticipant()?.username || 'Chat'
    },

    conversationPhoto() {
      if (this.conversation?.type === 'group') {
        return this.conversation.photoUrl || '/api/placeholder/40/40'
      }
      return this.getOtherParticipant()?.photoUrl || '/api/placeholder/40/40'
    },

    isGroup() {
      return this.conversation?.type === 'group'
    },
    filteredUsers() {
      if (!this.userSearchQuery) return this.users;
      const query = this.userSearchQuery.toLowerCase();
      
      // Filtra gli utenti che non sono già nel gruppo
      return this.users.filter(user => 
        user.username.toLowerCase().includes(query) &&
        !this.conversation.participants.some(p => p.id === user.id)
      );
    }
  },

  methods: {
    // Gestione dei messaggi base
    getOtherParticipant() {
      return this.conversation?.participants?.find(p => p.id !== this.currentUserId)
    },

    async fetchConversation() {
      try {
        const response = await this.$axios.get(`/conversations/${this.conversationID}`, {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        })
        this.conversation = response.data.conversation
        this.messages = response.data.messages
        this.scrollToBottom()
        this.markMessagesAsRead()
      } catch (error) {
        console.error('Errore nel recupero della conversazione:', error)
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

    async sendMessage() {
      if (!this.newMessage.trim()) return

      try {
        const messageData = {
          type: 'text',
          content: this.newMessage
        }

        await this.$axios.post(
          `/conversations/${this.conversationID}`,
          messageData,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        )
        
        this.newMessage = ''
        await this.fetchConversation()
      } catch (error) {
        console.error('Errore nell\'invio del messaggio:', error)
      }
    },

    // Menu contestuale dei messaggi
    openMessageMenu(message) {
      // Chiude il menu se è già aperto per lo stesso messaggio
      if (this.selectedMessage === message.id) {
        this.selectedMessage = null
      } else {
        this.selectedMessage = message.id
      }
      // Chiude altri modal aperti
      this.showCommentsModal = false
      this.showForwardModal = false
    },

    // Gestione commenti
    commentMessage(message) {
    // Apre direttamente il modale dei commenti per il messaggio selezionato
      this.showComments(message);
    },

    async showComments(message) {
      this.selectedMessageForComments = message;
      this.selectedMessage = null;
      this.isLoadingComments = true;
      this.commentsError = null;
      this.currentComments = []; // Reset array before loading
      this.showCommentsModal = true;

      try {
        const response = await this.$axios.get(
          `/conversations/${this.conversationID}/messages/${message.id}/comments`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        this.currentComments = response.data || []; // Assicurati che sia sempre un array
      } catch (error) {
        console.error('Errore nel recupero dei commenti:', error);
        this.commentsError = 'Errore nel caricamento dei commenti';
        this.currentComments = []; // Assicurati che sia un array vuoto in caso di errore
      } finally {
        this.isLoadingComments = false;
      }
    },

    async submitComment() {
      if (!this.newComment.trim() || !this.selectedMessageForComments) return;

      try {
        const response = await this.$axios.post(
          `/conversations/${this.conversationID}/messages/${this.selectedMessageForComments.id}/comments`,
          {
            content: this.newComment.trim()
          },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        if (response.data) {
          this.currentComments = [...this.currentComments, response.data];
          this.newComment = '';
        }
      } catch (error) {
        console.error('Errore nell\'invio del commento:', error);
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
        );
        // Rimuovi il commento dalla lista
        const index = this.currentComments.findIndex(c => c.id === comment.id);
        if (index !== -1) {
          this.currentComments.splice(index, 1);
        }
      } catch (error) {
        console.error('Errore nell\'eliminazione del commento:', error);
      }
    },

    // Gestione inoltro messaggi
    getConversationName(conversation) {
      if (conversation.type === 'group') {
        return conversation.name || 'Gruppo senza nome'
      }
      // Trova il partecipante che non è l'utente corrente
      const otherParticipant = conversation.participants?.find(
        p => p.id !== this.currentUserId
      )
      return otherParticipant ? otherParticipant.username : 'Chat'
    },
    async forwardMessage(message) {
      this.messageToForward = message;
      this.selectedMessage = null;
      try {
        const response = await this.$axios.get('/conversations/', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        });
        // Filtra la conversazione corrente dall'elenco
        this.forwardingConversations = response.data.filter(
          conv => conv.id !== this.conversationID
        );
        this.showForwardModal = true;
      } catch (error) {
        console.error('Errore nel recupero delle conversazioni:', error);
      }
    },

    async confirmForward(targetConversationId) {
      try {
        await this.$axios.post(
          `/conversations/${this.conversationID}/messages/${this.messageToForward.id}`,
          {
            id: targetConversationId
          },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        this.closeForwardModal();
        // Opzionale: mostra un messaggio di successo
        this.$emit('message-forwarded');
      } catch (error) {
        console.error('Errore nell\'inoltro del messaggio:', error);
      }
    },

    // Eliminazione messaggi
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
        this.selectedMessage = null
        await this.fetchConversation()
      } catch (error) {
        console.error('Errore nell\'eliminazione del messaggio:', error)
      }
    },

    // Utilità
    formatTime(timestamp) {
      if (!timestamp) return '';
      
      const date = new Date(timestamp);
      
      // Verifica se la data è valida
      if (isNaN(date.getTime())) {
        console.error('Invalid date:', timestamp);
        return '';
      }
      
      return date.toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit'
      });
    },

    getCurrentDate() {
      const now = new Date();
      return now.toLocaleDateString('it-IT', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
        year: 'numeric'
      });
    },

    scrollToBottom() {
      this.$nextTick(() => {
        const container = this.$refs.messagesContainer
        if (container) {
          container.scrollTop = container.scrollHeight
        }
      })
    },

    // Gestione modale
    closeCommentsModal() {
      this.showCommentsModal = false
      this.selectedMessageForComments = null
      this.currentComments = []
      this.newComment = ''
    },

    closeForwardModal() {
      this.showForwardModal = false
      this.messageToForward = null
      this.forwardingConversations = []
    },

    // Polling e lifecycle
    startPolling() {
      this.pollingInterval = setInterval(this.fetchConversation, 3000)
    },

    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval)
        this.pollingInterval = null
      }
    },

    // gestione gruppi
    toggleGroupMenu() {
      this.showGroupMenu = !this.showGroupMenu;
    },

    async openAddParticipantsModal() {
      this.showGroupMenu = false;
      this.showAddParticipantsModal = true;
      this.userSearchQuery = '';
      
      try {
        // Fetch users excluding current participants and current user
        const response = await this.$axios.get('/users/', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          },
          params: {
            // Optional: if backend supports filtering out existing participants
            exclude_conversation_id: this.conversationID
          }
        });
        
        // Filter out existing participants and current user
        this.users = response.data.filter(user => 
          user.id !== this.currentUserId && 
          !this.conversation.participants.some(p => p.id === user.id)
        );
      } catch (error) {
        console.error('Errore nel recupero degli utenti:', error);
        this.$emit('show-error', 'Impossibile recuperare gli utenti');
      }
    },

    async fetchUsers() {
      this.isLoadingUsers = true;
      try {
        const response = await this.$axios.get('/users/', {
          headers: {
            Authorization: `Bearer ${localStorage.getItem('authToken')}`
          }
        });
        this.users = response.data;
      } catch (error) {
        console.error('Errore nel recupero degli utenti:', error);
      } finally {
        this.isLoadingUsers = false;
      }
    },

    async addParticipant(user) {
      try {
        await this.$axios.post(
          `/conversations/${this.conversationID}/participants`,
          { userId: user.id },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        // Refresh conversation to update participant list
        await this.fetchConversation();
        
        // Close the modal
        this.showAddParticipantsModal = false;
        
        // Show success message
        this.$emit('show-success', `${user.username} è stato aggiunto al gruppo`);
      } catch (error) {
        console.error('Errore nell\'aggiunta del partecipante:', error);
        this.$emit('show-error', 'Impossibile aggiungere il partecipante');
      }
    },

    async removeParticipant(participant) {
      try {
        await this.$axios.delete(
          `/conversations/${this.conversationID}/participants`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`,
            },
            data: { userId: participant.id }
          }
        );
        
        // Refresh the conversation to update participant list
        await this.fetchConversation();
        
        this.$emit('show-success', `${participant.username} è stato rimosso dal gruppo`);
      } catch (error) {
        console.error('Errore nella rimozione del partecipante:', error);
        this.$emit('show-error', 'Impossibile rimuovere il partecipante');
      }
    },

    openEditGroupModal() {
      this.showGroupMenu = false;
      
      // Ensure we have a conversation and it's a group
      if (!this.conversation || this.conversation.type !== 'group') {
        console.error('Operazione non consentita');
        return;
      }

      this.editGroupName = this.conversation.name || '';
      this.editGroupPhoto = this.conversation.photoUrl || '';
      this.showEditGroupModal = true;
    },

    async saveGroupChanges() {
      // Validate input
      if (!this.editGroupName.trim()) {
        this.$emit('show-error', 'Il nome del gruppo non può essere vuoto');
        return;
      }

      try {
        // Update group name
        await this.$axios.put(
          `/conversations/${this.conversationID}/name`,
          { name: this.editGroupName.trim() },
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );

        // Update group photo if provided
        if (this.editGroupPhoto.trim()) {
          await this.$axios.put(
            `/conversations/${this.conversationID}/photo`,
            { photoUrl: this.editGroupPhoto.trim() },
            {
              headers: {
                Authorization: `Bearer ${localStorage.getItem('authToken')}`
              }
            }
          );
        }

        // Refresh conversation to show updates
        await this.fetchConversation();

        // Close the modal
        this.showEditGroupModal = false;

        // Show success message
        this.$emit('show-success', 'Gruppo aggiornato con successo');
      } catch (error) {
        console.error('Errore nell\'aggiornamento del gruppo:', error);
        
        // Determine the specific error message
        if (error.response && error.response.data) {
          this.$emit('show-error', error.response.data);
        } else {
          this.$emit('show-error', 'Impossibile aggiornare il gruppo');
        }
      }
    },

    async leaveGroup() {
      if (!confirm('Sei sicuro di voler uscire dal gruppo?')) return;

      try {
        await this.$axios.delete(
          `/conversations/${this.conversationID}/participants`,
          {
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          }
        );
        
        // Emetti un evento per informare il componente padre
        this.$emit('conversation-deleted', this.conversationID);
        
        // Reindirizza alla lista delle conversazioni
        this.$router.push('/conversations');
      } catch (error) {
        console.error('Errore durante l\'uscita dal gruppo:', error);
        this.$emit('show-error', 'Errore durante l\'uscita dal gruppo');
      }
    },

    closeAddParticipantsModal() {
      this.showAddParticipantsModal = false;
      this.userSearchQuery = '';
      this.users = [];
    },

    closeEditGroupModal() {
      this.showEditGroupModal = false;
      this.editGroupName = '';
      this.editGroupPhoto = '';
    }

  },

  async created() {
    await this.fetchConversation()
    this.startPolling()
  },

  beforeUnmount() {
    this.stopPolling()
  },

  watch: {
    conversationID: {
      immediate: true,
      handler(newId) {
        if (newId) {
          this.fetchConversation()
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

.chat-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 16px;
  background-color: #202c33;
  height: 60px;
  border-left: 1px solid #2a3942;
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}

.chat-header-info {
  display: flex;
  align-items: center;
  gap: 15px;
}

.chat-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.chat-details {
  display: flex;
  flex-direction: column;
}

.chat-name {
  color: #e9edef;
  font-size: 16px;
  font-weight: 500;
}

.chat-status {
  color: #8696a0;
  font-size: 13px;
}

.messages-container {
  flex: 1;
  overflow-y: auto;
  padding: 20px 5% 20px;
  position: relative;
  background-color: #0b141a;
  background-image: url("data:image/svg+xml,%3Csvg width='64' height='64' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath fill='%23192127' d='M8 16h48v32H8z' fill-opacity='.05'/%3E%3C/svg%3E");
}

.messages-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
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

.message-wrapper {
  display: flex;
  position: relative;
  max-width: 65%;
  margin: 2px 0;
}

.message-wrapper.sent {
  align-self: flex-end;
  flex-direction: row-reverse;
}

.message-wrapper.received {
  align-self: flex-start;
}

.message {
  padding: 8px 10px;
  border-radius: 7.5px;
  position: relative;
  min-width: 100px;
  cursor: context-menu;
  box-shadow: 0 1px 0.5px rgba(0, 0, 0, 0.13);
}

.sent .message {
  background-color: #005c4b;
  margin-right: 4px;
}

.received .message {
  background-color: #202c33;
  margin-left: 4px;
}

.message-sender {
  color: #53bdeb;
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 2px;
}

.message-content {
  position: relative;
  color: #e9edef;
  font-size: 14.2px;
  line-height: 19px;
}

.message-text {
  word-wrap: break-word;
  white-space: pre-wrap;
}

.message-media {
  max-width: 100%;
  border-radius: 4px;
  cursor: pointer;
}

.comments-badge {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #8696a0;
  font-size: 12px;
  margin-top: 4px;
}

.message-meta {
  display: flex;
  align-items: center;
  gap: 4px;
  float: right;
  margin-left: 8px;
  margin-bottom: -5px;
}

.message-time {
  color: #8696a0;
  font-size: 11px;
  margin-top: 4px;
}

.message-status {
  color: #8696a0;
  font-size: 11px;
  margin-top: 4px;
}

.message-status .read {
  color: #53bdeb;
}

.message-menu {
  position: absolute;
  background-color: #233138;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.4);
  z-index: 100;
  min-width: 200px;
  padding: 4px 0;
}

.menu-sent {
  right: 100%;
  margin-right: 8px;
  top: 0;
}

.message-wrapper.received .message-menu {
  left: 100%;
  margin-left: 8px;
  top: 0;
}

.menu-item {
  padding: 10px 16px;
  color: #e9edef;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 14px;
  transition: background-color 0.2s;
}

.menu-item:hover {
  background-color: #182229;
}

.menu-item i {
  width: 20px;
  text-align: center;
  font-size: 16px;
}

.menu-item.delete {
  color: #f15c6d;
}

.menu-item.delete i {
  color: #f15c6d;
}

.input-area {
  background-color: #202c33;
  padding: 10px 16px;
  position: relative;
  z-index: 2;
  flex-shrink: 0;
  border-left: 1px solid #2a3942;
}

.input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
  background-color: #2a3942;
  border-radius: 8px;
  padding: 9px 12px;
  position: relative;
}

.message-input {
  flex: 1;
  background: none;
  border: none;
  color: #d1d7db;
  font-size: 15px;
  outline: none;
  min-height: 20px;
  line-height: 20px;
  padding: 0;
}

.message-input::placeholder {
  color: #8696a0;
}

.send-button {
  background: none;
  border: none;
  color: #00a884;
  font-size: 20px;
  padding: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.send-button:disabled {
  color: #8696a0;
  cursor: default;
}

/* Modal styles */
.modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: #1f2c34;
  border-radius: 12px;
  width: 90%;
  max-width: 500px;
  max-height: 80vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.comments-modal {
  max-width: 450px;
}

.modal-header {
  padding: 16px 20px;
  border-bottom: 1px solid #2a3942;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background-color: #202c33;
}

.modal-header h3 {
  color: #e9edef;
  margin: 0;
  font-size: 16px;
  font-weight: 500;
}

.modal-input {
  width: 100%;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  color: #d1d7db;
  padding: 12px;
  margin-bottom: 16px;
  font-size: 15px;
  outline: none;
}

.modal-input::placeholder {
  color: #8696a0;
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
.modal-body {
  padding: 20px;
}
.modal-button {
  background-color: #00a884;
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  font-size: 14px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.modal-button:hover {
  background-color: #008069;
}

.modal-button:disabled {
  background-color: #2a3942;
  color: #8696a0;
  cursor: not-allowed;
}

/* Users list styles */
.users-list {
  max-height: 300px;
  overflow-y: auto;
}

.user-item {
  display: flex;
  align-items: center;
  padding: 12px;
  cursor: pointer;
  color: #d1d7db;
  border-radius: 8px;
  transition: background-color 0.2s;
}

.user-item:hover {
  background-color: #2a3942;
}

/* Search input specific styles */
.search-input {
  width: 100%;
  background-color: #2a3942;
  border: none;
  border-radius: 8px;
  color: #d1d7db;
  padding: 12px;
  font-size: 15px;
  outline: none;
  margin-bottom: 12px;
}

.search-input::placeholder {
  color: #8696a0;
}

/* Comments styles */
.comments-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  margin-bottom: 16px;
}

.no-comments {
  text-align: center;
  color: #8696a0;
  padding: 20px;
  font-size: 14px;
  background-color: rgba(134, 150, 160, 0.05);
  border-radius: 6px;
}

.comment {
  background-color: #202c33;
  padding: 12px 16px;
  border-radius: 8px;
  position: relative;
}

.comment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 4px;
}

.comment-sender {
  color: #00a884;
  font-size: 14px;
  font-weight: 500;
}

.comment-time {
  color: #8696a0;
  font-size: 12px;
}

.comment-content {
  color: #e9edef;
  font-size: 14px;
  line-height: 1.4;
  word-break: break-word;
}

.delete-comment {
  position: absolute;
  top: 8px;
  right: 8px;
  background: none;
  border: none;
  color: #8696a0;
  cursor: pointer;
  padding: 4px;
  border-radius: 50%;
  opacity: 0;
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
}

.comment:hover .delete-comment {
  opacity: 1;
}

.delete-comment:hover {
  background-color: #2a3942;
  color: #f15c6d;
}

.comments-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 20px;
  color: #8696a0;
}

.loading-spinner {
  width: 30px;
  height: 30px;
  border: 3px solid #2a3942;
  border-top-color: #00a884;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 10px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.comments-error {
  color: #f15c6d;
  background-color: rgba(241, 92, 109, 0.1);
  padding: 12px;
  border-radius: 6px;
  margin-bottom: 16px;
  text-align: center;
}


/* Forward modal styles */
.conversations-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.conversation-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-radius: 8px;
}

.conversation-item:hover {
  background-color: #2a3942;
}

.conversation-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
}

.conversation-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.conversation-name {
  color: #e9edef;
  font-size: 16px;
}

.conversation-type {
  color: #8696a0;
  font-size: 13px;
}

.no-conversations {
  text-align: center;
  color: #8696a0;
  padding: 20px;
  font-size: 14px;
}

/* Reply styles */
.reply-container {
  background-color: #1f2c34;
  padding: 8px 12px;
  margin-bottom: 8px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-left: 4px solid #00a884;
  border-radius: 4px;
}

.reply-content {
  flex: 1;
  min-width: 0;
}

.reply-sender {
  color: #00a884;
  font-size: 13px;
  font-weight: 500;
}

.reply-text {
  color: #d1d7db;
  font-size: 13px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.icon-button {
  background: none;
  border: none;
  color: #8696a0;
  padding: 4px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  transition: all 0.2s;
}

.icon-button:hover {
  background-color: #2a3942;
}

/* Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Custom scrollbar */
.messages-container::-webkit-scrollbar,
.modal-body::-webkit-scrollbar {
  width: 6px;
}

.messages-container::-webkit-scrollbar-track,
.modal-body::-webkit-scrollbar-track {
  background: transparent;
}

.messages-container::-webkit-scrollbar-thumb,
.modal-body::-webkit-scrollbar-thumb {
  background: rgba(134, 150, 160, 0.2);
  border-radius: 3px;
}

.messages-container::-webkit-scrollbar-thumb:hover,
.modal-body::-webkit-scrollbar-thumb:hover {
  background: rgba(134, 150, 160, 0.3);
}
.modal-body::-webkit-scrollbar {
  width: 6px;
}

.modal-body::-webkit-scrollbar-track {
  background: transparent;
}

.modal-body::-webkit-scrollbar-thumb {
  background: rgba(134, 150, 160, 0.2);
  border-radius: 3px;
}

.modal-body::-webkit-scrollbar-thumb:hover {
  background: rgba(134, 150, 160, 0.3);
}

/* New comment input styles */
.new-comment {
  display: flex;
  gap: 8px;
  background-color: #2a3942;
  border-radius: 8px;
  padding: 6px;
  margin-top: 16px;
}

.new-comment input {
  flex: 1;
  background: none;
  border: none;
  color: #e9edef;
  font-size: 14px;
  padding: 8px 12px;
  outline: none;
}

.new-comment input::placeholder {
  color: #8696a0;
}

.new-comment button {
  background: none;
  border: none;
  color: #00a884;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  transition: color 0.2s;
}

.new-comment button:disabled {
  color: #8696a0;
  cursor: default;
}

.new-comment button:not(:disabled):hover {
  color: #00c29d;
}

/* Quote styles */
.quoted-message {
  background-color: #ffffff1a;
  border-left: 4px solid #53bdeb;
  border-radius: 4px;
  margin-bottom: 8px;
  padding: 5px 10px;
}

.quoted-sender {
  color: #53bdeb;
  font-size: 13px;
  font-weight: 500;
}

.quoted-content {
  color: #e9edef;
  font-size: 13px;
  opacity: 0.7;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* Message status indicators */
.message-status i {
  font-size: 14px;
}

.message-status i.read {
  color: #53bdeb;
}

/* Message context menu animations */
.message-menu {
  animation: fadeIn 0.2s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Modal animations */
.modal {
  animation: modalFadeIn 0.3s ease-out;
}

.modal-content {
  animation: modalSlideIn 0.3s ease-out;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes modalSlideIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Responsive design */
@media (max-width: 768px) {
  .message-wrapper {
    max-width: 85%;
  }

  .modal-content {
    width: 95%;
    max-height: 90vh;
  }

  .chat-header {
    padding: 8px 12px;
  }

  .input-area {
    padding: 8px 12px;
  }

  .message {
    font-size: 14px;
    padding: 6px 8px;
  }

  .message-menu {
    min-width: 180px;
  }

  .menu-item {
    padding: 8px 12px;
    font-size: 13px;
  }
}

/* Dark theme enhancements */
.message-wrapper.sent .message::before {
  content: "";
  position: absolute;
  top: 0;
  right: -8px;
  width: 8px;
  height: 13px;
  background: linear-gradient(
    to top left,
    transparent 50%,
    #005c4b 50%
  );
}

.message-wrapper.received .message::before {
  content: "";
  position: absolute;
  top: 0;
  left: -8px;
  width: 8px;
  height: 13px;
  background: linear-gradient(
    to top right,
    transparent 50%,
    #202c33 50%
  );
}

/* Input area enhancements */
.input-wrapper:focus-within {
  box-shadow: 0 0 0 2px #00a884;
}

.new-comment:focus-within {
  box-shadow: 0 0 0 2px #00a884;
}

/* Emoji styles */
.emoji {
  font-size: 1.2em;
  line-height: 1;
  vertical-align: -0.1em;
}

/* Loading states */
.loading {
  opacity: 0.7;
  pointer-events: none;
}

.loading-spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid #8696a0;
  border-top-color: transparent;
  border-radius: 50%;
  animation: spinner 0.8s linear infinite;
}

@keyframes spinner {
  to {
    transform: rotate(360deg);
  }
}

/* Error states */
.error-message {
  color: #f15c6d;
  font-size: 14px;
  padding: 8px 12px;
  background-color: rgba(241, 92, 109, 0.1);
  border-radius: 4px;
  margin: 8px 0;
}

/* Success states */
.success-message {
  color: #00a884;
  font-size: 14px;
  padding: 8px 12px;
  background-color: rgba(0, 168, 132, 0.1);
  border-radius: 4px;
  margin: 8px 0;
}

/* Message selection styles */
.message.selected {
  box-shadow: 0 0 0 2px #00a884;
}

/* Unread message indicator */
.unread-indicator {
  background-color: #00a884;
  width: 6px;
  height: 6px;
  border-radius: 50%;
  margin-left: 4px;
  display: inline-block;
}
/*  Group styles  */
.chat-header-actions {
  position: relative;
}

.group-menu {
  position: absolute;
  top: 100%;
  right: 0;
  background-color: #233138;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.2);
  z-index: 1000;
  min-width: 220px;
  padding: 8px 0;
  margin-top: 8px;
}

.participants-modal .modal-body {
  padding: 0;
}

.participants-list {
  max-height: 400px;
  overflow-y: auto;
}

.participant-item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid #2a3942;
}

.participant-item:last-child {
  border-bottom: none;
}

.participant-avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  margin-right: 12px;
  object-fit: cover;
}

.participant-info {
  flex: 1;
}

.participant-name {
  color: #e9edef;
  font-size: 16px;
}

.participant-you {
  color: #8696a0;
  font-size: 14px;
  margin-left: 8px;
}

.remove-participant-btn {
  background: none;
  border: none;
  color: #f15c6d;
  cursor: pointer;
  padding: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.remove-participant-btn:hover {
  background-color: rgba(241, 92, 109, 0.1);
}
</style>