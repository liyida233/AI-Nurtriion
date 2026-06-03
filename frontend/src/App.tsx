import {
  Activity,
  Apple,
  Bot,
  Dumbbell,
  FileText,
  Flag,
  LogOut,
  RefreshCw,
  Scale,
  UserRound
} from "lucide-react";
import { FormEvent, useEffect, useMemo, useState } from "react";
import { api } from "./api";

type AuthResponse = {
  token: string;
  user: { id: string; name: string; email: string; role: string };
};

type DashboardSummary = {
  period: string;
  workoutSessions: number;
  trainingVolume: number;
  workoutConsistency: number;
  caloriesIn: number;
  protein: number;
  carbohydrates: number;
  fat: number;
  estimatedBmr: number;
  estimatedTdee: number;
  calorieBalance: number;
  latestWeightKg: number;
  weightTrend: string;
  activeGoals: number;
  goalAdherence: number;
};

type Recommendation = {
  id: string;
  type: string;
  content: string;
  createdAt: string;
};

type ProgressReport = {
  id: string;
  periodType: string;
  generatedAt: string;
  summary: string;
};

type Exercise = {
  id: string;
  name: string;
  muscleGroup: string;
};

type FoodItem = {
  id: string;
  name: string;
  servingSize: string;
};

const navItems = [
  { key: "dashboard", label: "Dashboard", icon: Activity },
  { key: "profile", label: "Profile", icon: UserRound },
  { key: "workout", label: "Workout", icon: Dumbbell },
  { key: "nutrition", label: "Nutrition", icon: Apple },
  { key: "body", label: "Body", icon: Scale },
  { key: "goals", label: "Goals", icon: Flag },
  { key: "reports", label: "Reports", icon: FileText },
  { key: "ai", label: "AI", icon: Bot }
] as const;

type NavKey = (typeof navItems)[number]["key"];

export default function App() {
  const [token, setToken] = useState(() => localStorage.getItem("token") ?? "");
  const [user, setUser] = useState<AuthResponse["user"] | null>(() => {
    const cached = localStorage.getItem("user");
    return cached ? JSON.parse(cached) : null;
  });
  const [active, setActive] = useState<NavKey>("dashboard");
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [recommendations, setRecommendations] = useState<Recommendation[]>([]);
  const [reports, setReports] = useState<ProgressReport[]>([]);
  const [message, setMessage] = useState("");

  const isAuthed = Boolean(token);

  async function loadDashboard() {
    if (!token) return;
    const result = await api<DashboardSummary>("/dashboard/summary", {}, token);
    if (result.error) {
      setMessage(result.error);
      return;
    }
    setSummary(result.data ?? null);
  }

  async function loadRecommendations() {
    if (!token) return;
    const result = await api<Recommendation[]>("/recommendations", {}, token);
    if (!result.error) {
      setRecommendations(result.data ?? []);
    }
  }

  async function loadReports() {
    if (!token) return;
    const result = await api<ProgressReport[]>("/reports", {}, token);
    if (!result.error) {
      setReports(result.data ?? []);
    }
  }

  useEffect(() => {
    if (isAuthed) {
      loadDashboard();
      loadRecommendations();
      loadReports();
    }
  }, [isAuthed]);

  function saveAuth(auth: AuthResponse) {
    localStorage.setItem("token", auth.token);
    localStorage.setItem("user", JSON.stringify(auth.user));
    setToken(auth.token);
    setUser(auth.user);
    setMessage("");
  }

  function logout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setToken("");
    setUser(null);
    setSummary(null);
  }

  if (!isAuthed) {
    return <AuthScreen onAuth={saveAuth} />;
  }

  const ActiveIcon = navItems.find((item) => item.key === active)?.icon ?? Activity;

  return (
    <main className="app-shell">
      <aside className="sidebar">
        <div className="brand">
          <div className="brand-mark">AI</div>
          <div>
            <strong>Nutrition</strong>
            <span>Decision Support</span>
          </div>
        </div>
        <nav>
          {navItems.map((item) => {
            const Icon = item.icon;
            return (
              <button
                className={active === item.key ? "nav-item active" : "nav-item"}
                key={item.key}
                onClick={() => setActive(item.key)}
                title={item.label}
              >
                <Icon size={18} />
                <span>{item.label}</span>
              </button>
            );
          })}
        </nav>
        <button className="nav-item logout" onClick={logout} title="Log out">
          <LogOut size={18} />
          <span>Log out</span>
        </button>
      </aside>

      <section className="workspace">
        <header className="topbar">
          <div>
            <p>{user?.name ?? "User"}</p>
            <h1>
              <ActiveIcon size={24} />
              {navItems.find((item) => item.key === active)?.label}
            </h1>
          </div>
          <button className="icon-button" onClick={loadDashboard} title="Refresh dashboard">
            <RefreshCw size={18} />
          </button>
        </header>

        {message && <div className="notice">{message}</div>}

        {active === "dashboard" && <Dashboard summary={summary} />}
        {active === "profile" && <ProfileForm token={token} onSaved={loadDashboard} setMessage={setMessage} />}
        {active === "workout" && <WorkoutPanel token={token} onSaved={loadDashboard} setMessage={setMessage} />}
        {active === "nutrition" && <NutritionPanel token={token} onSaved={loadDashboard} setMessage={setMessage} />}
        {active === "body" && <BodyPanel token={token} onSaved={loadDashboard} setMessage={setMessage} />}
        {active === "goals" && <GoalPanel token={token} onSaved={loadDashboard} setMessage={setMessage} />}
        {active === "reports" && (
          <ReportPanel token={token} reports={reports} reload={loadReports} setMessage={setMessage} />
        )}
        {active === "ai" && (
          <AIPanel
            token={token}
            recommendations={recommendations}
            reload={loadRecommendations}
            setMessage={setMessage}
          />
        )}
      </section>
    </main>
  );
}

function AuthScreen({ onAuth }: { onAuth: (auth: AuthResponse) => void }) {
  const [mode, setMode] = useState<"login" | "register">("login");
  const [form, setForm] = useState({ name: "", email: "", password: "" });
  const [error, setError] = useState("");

  async function submit(event: FormEvent) {
    event.preventDefault();
    const result = await api<AuthResponse>(`/auth/${mode}`, {
      method: "POST",
      body: JSON.stringify(mode === "login" ? { email: form.email, password: form.password } : form)
    });
    if (result.error || !result.data) {
      setError(result.error ?? "Authentication failed");
      return;
    }
    onAuth(result.data);
  }

  return (
    <main className="auth-screen">
      <section className="auth-panel">
        <div className="brand large">
          <div className="brand-mark">AI</div>
          <div>
            <strong>AI Nutrition</strong>
            <span>Fitness and nutrition tracking</span>
          </div>
        </div>
        <form onSubmit={submit} className="form-grid">
          {mode === "register" && (
            <label>
              Name
              <input value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
            </label>
          )}
          <label>
            Email
            <input
              type="email"
              value={form.email}
              onChange={(e) => setForm({ ...form, email: e.target.value })}
              required
            />
          </label>
          <label>
            Password
            <input
              type="password"
              value={form.password}
              onChange={(e) => setForm({ ...form, password: e.target.value })}
              minLength={8}
              required
            />
          </label>
          {error && <div className="notice danger">{error}</div>}
          <button className="primary" type="submit">
            {mode === "login" ? "Log in" : "Create account"}
          </button>
          <button className="text-button" type="button" onClick={() => setMode(mode === "login" ? "register" : "login")}>
            {mode === "login" ? "Create an account" : "Use existing account"}
          </button>
        </form>
      </section>
    </main>
  );
}

function Dashboard({ summary }: { summary: DashboardSummary | null }) {
  const cards = useMemo(
    () => [
      ["Workout Sessions", summary?.workoutSessions ?? 0],
      ["Training Volume", `${Math.round(summary?.trainingVolume ?? 0)} kg`],
      ["Calories In", `${Math.round(summary?.caloriesIn ?? 0)} kcal`],
      ["Protein", `${Math.round(summary?.protein ?? 0)} g`],
      ["Calorie Balance", `${Math.round(summary?.calorieBalance ?? 0)} kcal`],
      ["Weight Trend", summary?.weightTrend ?? "No data"],
      ["Active Goals", summary?.activeGoals ?? 0],
      ["Goal Adherence", `${Math.round(summary?.goalAdherence ?? 0)}%`]
    ],
    [summary]
  );

  return (
    <section className="content-grid">
      {cards.map(([label, value]) => (
        <article className="metric" key={label}>
          <span>{label}</span>
          <strong>{value}</strong>
        </article>
      ))}
    </section>
  );
}

function ProfileForm({ token, onSaved, setMessage }: PanelProps) {
  return (
    <SimpleForm
      token={token}
      path="/profile"
      method="PUT"
      onSaved={onSaved}
      setMessage={setMessage}
      fields={[
        ["age", "Age", "number", "22"],
        ["gender", "Gender", "text", "male"],
        ["heightCm", "Height cm", "number", "175"],
        ["weightKg", "Weight kg", "number", "70"],
        ["activityLevel", "Activity level", "text", "moderate"],
        ["primaryGoal", "Primary goal", "text", "fat_loss"]
      ]}
    />
  );
}

function WorkoutPanel({ token, onSaved, setMessage }: PanelProps) {
  const [exercises, setExercises] = useState<Exercise[]>([]);

  useEffect(() => {
    api<Exercise[]>("/workouts/exercises", {}, token).then((result) => {
      if (!result.error) {
        setExercises(result.data ?? []);
      }
    });
  }, [token]);

  return (
    <SimpleForm
      token={token}
      path="/workouts"
      method="POST"
      onSaved={onSaved}
      setMessage={setMessage}
      transform={(values) => ({
        workoutDate: values.workoutDate,
        category: values.category,
        durationMin: Number(values.durationMin),
        notes: values.notes,
        entries: [
          {
            exerciseId: values.exerciseId,
            sets: Number(values.sets),
            reps: Number(values.reps),
            weightKg: Number(values.weightKg),
            restSec: Number(values.restSec)
          }
        ]
      })}
      fields={[
        ["workoutDate", "Workout date", "date", new Date().toISOString().slice(0, 10)],
        ["category", "Category", "text", "strength"],
        ["durationMin", "Duration min", "number", "45"],
        [
          "exerciseId",
          "Exercise",
          "select",
          exercises[0]?.id ?? "",
          exercises.map((exercise) => ({
            value: exercise.id,
            label: `${exercise.name} · ${exercise.muscleGroup}`
          }))
        ],
        ["sets", "Sets", "number", "3"],
        ["reps", "Reps", "number", "10"],
        ["weightKg", "Weight kg", "number", "40"],
        ["restSec", "Rest sec", "number", "60"],
        ["notes", "Notes", "text", ""]
      ]}
    />
  );
}

function NutritionPanel({ token, onSaved, setMessage }: PanelProps) {
  const [foods, setFoods] = useState<FoodItem[]>([]);

  useEffect(() => {
    api<FoodItem[]>("/nutrition/foods", {}, token).then((result) => {
      if (!result.error) {
        setFoods(result.data ?? []);
      }
    });
  }, [token]);

  return (
    <SimpleForm
      token={token}
      path="/nutrition/meals"
      method="POST"
      onSaved={onSaved}
      setMessage={setMessage}
      transform={(values) => ({
        foodItemId: values.foodItemId,
        mealType: values.mealType,
        quantity: Number(values.quantity),
        mealTime: new Date(values.mealTime).toISOString()
      })}
      fields={[
        [
          "foodItemId",
          "Food",
          "select",
          foods[0]?.id ?? "",
          foods.map((food) => ({
            value: food.id,
            label: `${food.name} · ${food.servingSize}`
          }))
        ],
        ["mealType", "Meal type", "text", "lunch"],
        ["quantity", "Quantity", "number", "1"],
        ["mealTime", "Meal time", "datetime-local", ""]
      ]}
    />
  );
}

function BodyPanel({ token, onSaved, setMessage }: PanelProps) {
  return (
    <SimpleForm
      token={token}
      path="/body-records"
      method="POST"
      onSaved={onSaved}
      setMessage={setMessage}
      transform={(values) => ({ ...values, weightKg: Number(values.weightKg) })}
      fields={[
        ["recordDate", "Record date", "date", new Date().toISOString().slice(0, 10)],
        ["weightKg", "Weight kg", "number", "70"],
        ["note", "Note", "text", ""]
      ]}
    />
  );
}

function GoalPanel({ token, onSaved, setMessage }: PanelProps) {
  return (
    <SimpleForm
      token={token}
      path="/goals"
      method="POST"
      onSaved={onSaved}
      setMessage={setMessage}
      transform={(values) => ({ ...values, targetValue: Number(values.targetValue) })}
      fields={[
        ["goalType", "Goal type", "text", "workout_frequency"],
        ["targetMetric", "Target metric", "text", "sessions_per_week"],
        ["targetValue", "Target value", "number", "3"],
        ["deadline", "Deadline", "date", ""],
        ["priority", "Priority", "text", "high"]
      ]}
    />
  );
}

function AIPanel({
  token,
  recommendations,
  reload,
  setMessage
}: {
  token: string;
  recommendations: Recommendation[];
  reload: () => void;
  setMessage: (message: string) => void;
}) {
  async function generate(type: string) {
    const result = await api<Recommendation>(
      "/recommendations/generate",
      { method: "POST", body: JSON.stringify({ type }) },
      token
    );
    if (result.error) {
      setMessage(result.error);
      return;
    }
    setMessage("Recommendation generated");
    reload();
  }

  return (
    <section className="panel-stack">
      <div className="action-row">
        <button className="primary" onClick={() => generate("weekly")}>
          Generate weekly feedback
        </button>
        <button onClick={() => generate("workout")}>Workout plan</button>
        <button onClick={() => generate("meal")}>Meal idea</button>
      </div>
      {recommendations.map((item) => (
        <article className="recommendation" key={item.id}>
          <span>{item.type}</span>
          <pre>{item.content}</pre>
        </article>
      ))}
    </section>
  );
}

function ReportPanel({
  token,
  reports,
  reload,
  setMessage
}: {
  token: string;
  reports: ProgressReport[];
  reload: () => void;
  setMessage: (message: string) => void;
}) {
  async function generate(periodType: "weekly" | "monthly") {
    const result = await api<{ report: ProgressReport; metrics: DashboardSummary }>(
      "/reports/generate",
      { method: "POST", body: JSON.stringify({ periodType }) },
      token
    );
    if (result.error) {
      setMessage(result.error);
      return;
    }
    setMessage(`${periodType} report generated`);
    reload();
  }

  return (
    <section className="panel-stack">
      <div className="action-row">
        <button className="primary" onClick={() => generate("weekly")}>
          Generate weekly report
        </button>
        <button onClick={() => generate("monthly")}>Generate monthly report</button>
      </div>
      {reports.length === 0 && <div className="notice">No reports yet. Generate one from current analytics.</div>}
      {reports.map((report) => (
        <article className="recommendation" key={report.id}>
          <span>
            {report.periodType} · {new Date(report.generatedAt).toLocaleString()}
          </span>
          <pre>{report.summary}</pre>
        </article>
      ))}
    </section>
  );
}

type PanelProps = {
  token: string;
  onSaved: () => void;
  setMessage: (message: string) => void;
};

type SelectOption = {
  value: string;
  label: string;
};

type Field = [key: string, label: string, type: string, placeholder: string, options?: SelectOption[]];

function SimpleForm({
  token,
  path,
  method,
  fields,
  transform,
  onSaved,
  setMessage
}: PanelProps & {
  path: string;
  method: "POST" | "PUT";
  fields: Field[];
  transform?: (values: Record<string, string>) => unknown;
}) {
  const initial = Object.fromEntries(fields.map(([key, , , value]) => [key, value]));
  const [values, setValues] = useState<Record<string, string>>(initial);

  useEffect(() => {
    setValues((current) => {
      const next = { ...current };
      let changed = false;
      for (const [key, , , value] of fields) {
        if (!next[key] && value) {
          next[key] = value;
          changed = true;
        }
      }
      return changed ? next : current;
    });
  }, [fields]);

  async function submit(event: FormEvent) {
    event.preventDefault();
    const payload = transform ? transform(values) : coerceNumbers(values);
    const result = await api(path, { method, body: JSON.stringify(payload) }, token);
    if (result.error) {
      setMessage(result.error);
      return;
    }
    setMessage("Saved");
    onSaved();
  }

  return (
    <form className="data-form" onSubmit={submit}>
      {fields.map(([key, label, type, placeholder, options]) => (
        <label key={key}>
          {label}
          {type === "select" ? (
            <select
              value={values[key] ?? ""}
              onChange={(event) => setValues({ ...values, [key]: event.target.value })}
              required
            >
              <option value="">Select</option>
              {(options ?? []).map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          ) : (
            <input
              type={type}
              value={values[key] ?? ""}
              placeholder={placeholder}
              onChange={(event) => setValues({ ...values, [key]: event.target.value })}
              required={key !== "notes" && key !== "note"}
            />
          )}
        </label>
      ))}
      <button className="primary" type="submit">
        Save
      </button>
    </form>
  );
}

function coerceNumbers(values: Record<string, string>) {
  return Object.fromEntries(
    Object.entries(values).map(([key, value]) => {
      if (value !== "" && !Number.isNaN(Number(value))) {
        return [key, Number(value)];
      }
      return [key, value];
    })
  );
}
