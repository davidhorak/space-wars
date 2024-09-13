import { create } from 'zustand';

type Status = 'initialization' | 'ready' | 'error';

type State = {
    status: Status;
    error?: Error;
}

type Action = {
    setStatus: (status: Status) => void;
    setError: (error?: Error) => void;
}

const useAppStore = create<State & Action>()((set) => ({
    status: 'initialization',
    setStatus: (status: Status) => set({ status }),
    setError: (error?: Error) => set({ error }),
}))

export default useAppStore;
