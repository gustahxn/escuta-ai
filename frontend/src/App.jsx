import { useState, useCallback } from "react";
import { Youtube, Music2, Play } from "lucide-react";

function App() {
  const [query, setQuery] = useState("");
  const [searchResults, setSearchResults] = useState([]);
  const [recommendations, setRecommendations] = useState([]);
  const [selectedTrack, setSelectedTrack] = useState(null);
  const [loading, setLoading] = useState(false);

  const search = useCallback(async () => {
    if (!query.trim()) return;
    setLoading(true);
    setRecommendations([]);
    setSelectedTrack(null);
    try {
      const res = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
      const data = await res.json();
      setSearchResults(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  }, [query]);

  const selectTrack = async (track) => {
    setLoading(true);
    setSelectedTrack(track);
    setSearchResults([]);
    try {
      const res = await fetch(
        `/api/recommend?artist=${encodeURIComponent(
          track.artist
        )}&track=${encodeURIComponent(track.name)}`
      );
      const data = await res.json();
      setRecommendations(Array.isArray(data) ? data : []);
    } catch (err) {
      console.error(err);
      setRecommendations([]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-white text-black font-['Inter',sans-serif]">
      <nav className="border-b border-gray-200 py-6 sticky top-0 bg-white/95 backdrop-blur-xl z-20">
        <div className="max-w-5xl mx-auto px-6 flex flex-col md:flex-row items-center justify-between gap-6">
          <h1
            className="text-2xl font-black tracking-tighter cursor-pointer hover:opacity-60 transition-opacity"
            onClick={() => window.location.reload()}
          >
            ESCUTA AÍ{" "}
            <span className="font-light text-gray-600">
              / music auto matcher!
            </span>
          </h1>

          <div className="flex w-full md:w-auto bg-gray-50 px-5 py-3.5 focus-within:bg-white border border-gray-300 focus-within:border-gray-900 transition-all">
            <input
              className="bg-transparent outline-none px-2 w-full md:w-80 text-base font-medium placeholder:text-gray-400"
              placeholder="Pesquise uma música..."
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && search()}
            />
            <button
              onClick={search}
              className="text-xs font-black uppercase tracking-wider px-3 hover:opacity-60 transition-opacity"
            >
              Buscar
            </button>
          </div>
        </div>
      </nav>

      <main className="max-w-5xl mx-auto px-6 py-16">
        {loading && (
          <div className="text-center py-32 font-black text-sm tracking-[0.3em] animate-pulse text-gray-600">
            CARREGANDO...
          </div>
        )}

        {searchResults.length > 0 && (
          <div className="max-w-3xl mx-auto space-y-3">
            <p className="text-xs font-bold text-gray-600 uppercase tracking-wider mb-8">
              Selecione a versão correta:
            </p>
            {searchResults.map((t, i) => (
              <div
                key={i}
                onClick={() => selectTrack(t)}
                className="group flex justify-between items-center p-8 bg-black text-white border border-gray-800 hover:bg-zinc-900 hover:border-gray-700 cursor-pointer transition-all"
              >
                <div>
                  <h3 className="font-bold text-xl mb-1">{t.name}</h3>
                  <p className="text-gray-300 font-medium text-base">
                    {t.artist}
                  </p>
                </div>
                <div className="w-10 h-10 border border-gray-600 flex items-center justify-center group-hover:bg-white group-hover:text-black group-hover:border-white transition-all text-sm font-bold">
                  →
                </div>
              </div>
            ))}
          </div>
        )}

        {selectedTrack && !loading && (
          <div className="animate-in fade-in duration-700">
            <header className="mb-20 text-center">
              <p className="text-xs font-bold text-gray-600 uppercase tracking-wider mb-6">
                Análise baseada em:
              </p>
              <h2 className="text-7xl font-black tracking-tighter mb-6 leading-tight">
                {selectedTrack.name}
              </h2>
              <p className="text-3xl font-light text-gray-600">
                {selectedTrack.artist}
              </p>
            </header>

            {recommendations.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-x-12 gap-y-10">
                {recommendations.map((track, i) => (
                  <div
                    key={i}
                    className="group bg-black text-white border border-gray-800 p-8 hover:bg-zinc-900 hover:border-gray-700 hover:shadow-2xl transition-all"
                  >
                    <div className="flex justify-between items-start mb-4">
                      <span className="text-xs font-bold text-gray-300 uppercase tracking-wider">
                        Sugestão #{i + 1}
                      </span>
                    </div>
                    <h4 className="font-bold text-2xl leading-tight mb-2">
                      {track.name}
                    </h4>
                    <p className="text-gray-300 font-medium text-lg mb-6">
                      {track.artist?.name}
                    </p>
                    <div className="flex gap-2">
                      <a
                        href={`https://www.youtube.com/results?search_query=${encodeURIComponent(
                          track.artist?.name + " " + track.name
                        )}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex-1 flex items-center justify-center gap-2 text-xs font-black bg-white text-black py-3 hover:bg-gray-200 transition-colors no-underline"
                      >
                        <Youtube size={16} />{" "}
                        <span className="hidden sm:inline">YOUTUBE</span>
                      </a>
                      <a
                        href={`https://open.spotify.com/search/${encodeURIComponent(
                          track.artist?.name + " " + track.name
                        )}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex-1 flex items-center justify-center gap-2 text-xs font-black bg-white text-black py-3 hover:bg-gray-200 transition-colors no-underline"
                      >
                        <Music2 size={16} />{" "}
                        <span className="hidden sm:inline">SPOTIFY</span>
                      </a>
                      <a
                        href={`https://music.apple.com/br/search?term=${encodeURIComponent(
                          track.artist?.name + " " + track.name
                        )}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="flex-1 flex items-center justify-center gap-2 text-xs font-black bg-white text-black py-3 hover:bg-gray-200 transition-colors no-underline"
                      >
                        <Play size={16} />{" "}
                        <span className="hidden sm:inline">APPLE</span>
                      </a>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-center py-20 border-2 border-dashed border-gray-200">
                <p className="text-lg font-medium text-gray-600 uppercase tracking-widest">
                  Nenhuma recomendação encontrada.
                </p>
              </div>
            )}
          </div>
        )}

        {!loading && !selectedTrack && searchResults.length === 0 && (
          <div className="text-center py-40">
            <h2 className="text-8xl font-black tracking-tighter">
              MERGULHA NO SOM
            </h2>
          </div>
        )}
      </main>
    </div>
  );
}

export default App;
