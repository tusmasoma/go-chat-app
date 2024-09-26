export interface Message {
    id: string;
    content: string;
    role: 'user' | 'system' | 'assistant'; // roleに変更
}