import {
  Activity,
  Apple,
  Bot,
  CheckCircle,
  Download,
  Dumbbell,
  FileText,
  Flag,
  LogOut,
  RefreshCw,
  Scale,
  Shield,
  Trash2,
  UserRound
} from "lucide-react";
import { FormEvent, ReactNode, useEffect, useMemo, useState } from "react";
import { api, downloadFile } from "./api";

type User = {
  id: string;
  name: string;
  email: string;
  role: string;
  createdAt?: string;
};

type AuthResponse = {
  token: string;
  user: User;
};

type UserProfile = {
  id?: string;
  age: number;
  gender: string;
  heightCm: number;
  weightKg: number;
  activityLevel: string;
  primaryGoal: string;
};

type DashboardSummary = {
  period: string;
  generatedAt: string;
  startDate: string;
  endDate: string;
  days: number;
  workoutSessions: number;
  trainingVolume: number;
  workoutConsistency: number;
  progressiveOverloadStatus: string;
  muscleGroupDistribution: Record<string, number>;
  muscleGroupWarnings: string[];
  caloriesIn: number;
  protein: number;
  carbohydrates: number;
  fat: number;
  mealCount: number;
  mealQualityScore: number;
  proteinRatio: number;
  carbohydrateRatio: number;
  fatRatio: number;
  nutritionGaps: string[];
  estimatedBmr: number;
  estimatedTdee: number;
  calorieBalance: number;
  calorieStatus: string;
  latestWeightKg: number;
  weightMovingAverage7DayKg: number;
  weightTrend: string;
  activeGoals: number;
  completedMilestones: number;
  totalMilestones: number;
  goalAdherence: number;
};

type Exercise = {
  id: string;
  name: string;
  category: string;
  muscleGroup: string;
  equipment: string;
  intensityLevel: string;
};

type WorkoutEntry = {
  id: string;
  exerciseId: string;
  exercise?: Exercise;
  sets: number;
  reps: number;
  weightKg: number;
  restSec: number;
};

type WorkoutSession = {
  id: string;
  workoutDate: string;
  category: string;
  durationMin: number;
  notes: string;
  entries: WorkoutEntry[];
};

type FoodItem = {
  id: string;
  name: string;
  servingSize: string;
  calories: number;
  protein: number;
  carbohydrates: number;
  fat: number;
  sugar: number;
  sodium: number;
};

type MealLog = {
  id: string;
  foodItemId: string;
  foodItem?: FoodItem;
  mealType: string;
  quantity: number;
  mealTime: string;
  totalCalories: number;
  totalProtein: number;
  totalCarbs: number;
  totalFat: number;
};

type BodyRecord = {
  id: string;
  recordDate: string;
  weightKg: number;
  note: string;
};

type GoalMilestone = {
  id: string;
  title: string;
  targetValue: number;
  dueDate: string;
  completed: boolean;
};

type Goal = {
  id: string;
  goalType: string;
  targetMetric: string;
  targetValue: number;
  deadline: string;
  priority: string;
  status: string;
  milestones?: GoalMilestone[];
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
  fileUrl: string;
};

type NavItem = {
  key: NavKey;
  label: string;
  icon: typeof Activity;
  adminOnly?: boolean;
};

type NavKey = "dashboard" | "profile" | "workout" | "nutrition" | "body" | "goals" | "reports" | "ai" | "admin";

type PanelProps = {
  token: string;
  onChanged: () => void;
  setMessage: (message: string) => void;
};

const navItems: NavItem[] = [
  { key: "dashboard", label: "Dashboard", icon: Activity },
  { key: "profile", label: "Profile", icon: UserRound },
  { key: "workout", label: "Workout", icon: Dumbbell },
  { key: "nutrition", label: "Nutrition", icon: Apple },
  { key: "body", label: "Body", icon: Scale },
  { key: "goals", label: "Goals", icon: Flag },
  { key: "reports", label: "Reports", icon: FileText },
  { key: "ai", label: "AI", icon: Bot },
  { key: "admin", label: "Admin", icon: Shield, adminOnly: true }
];

export default function App() {
  const [token, setToken] = useState(() => localStorage.getItem("token") ?? "");
  const [user, setUser] = useState<User | null>(() => {
    const cached = localStorage.getItem("user");
    return cached ? JSON.parse(cached) : null;
  });
  const [active, setActive] = useState<NavKey>("dashboard");
  const [summary, setSummary] = useState<DashboardSummary | null>(null);
  const [analyticsRange, setAnalyticsRange] = useState({
    period: "weekly",
    startDate: today(),
    endDate: today()
  });
  const [message, setMessage] = useState("");

  const visibleNav = navItems.filter((item) => !item.adminOnly || user?.role === "admin");
  const activeItem = visibleNav.find((item) => item.key === active) ?? visibleNav[0];
  const ActiveIcon = activeItem.icon;

  async function loadDashboard() {
    if (!token) return;
    const query = dashboardQuery(analyticsRange);
    const result = await api<DashboardSummary>(`/dashboard/summary${query}`, {}, token);
    if (result.error) {
      setMessage(result.error);
      return;
    }
    setSummary(result.data ?? null);
  }

  async function refreshAll() {
    await loadDashboard();
  }

  useEffect(() => {
    if (token) {
      loadDashboard();
      api<User>("/auth/me", {}, token).then((result) => {
        if (result.data) {
          localStorage.setItem("user", JSON.stringify(result.data));
          setUser(result.data);
        }
      });
    }
  }, [token]);

  function saveAuth(auth: AuthResponse) {
    localStorage.setItem("token", auth.token);
    localStorage.setItem("user", JSON.stringify(auth.user));
    setToken(auth.token);
    setUser(auth.user);
    setMessage("");
  }

  async function logout() {
    if (token) {
      await api("/auth/logout", { method: "POST" }, token);
    }
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setToken("");
    setUser(null);
    setSummary(null);
    setActive("dashboard");
  }

  if (!token) {
    return <AuthScreen onAuth={saveAuth} />;
  }

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
          {visibleNav.map((item) => {
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
            <p>
              {user?.name ?? "User"} / {user?.role ?? "user"}
            </p>
            <h1>
              <ActiveIcon size={24} />
              {activeItem.label}
            </h1>
          </div>
          <button className="icon-button" onClick={refreshAll} title="Refresh">
            <RefreshCw size={18} />
          </button>
        </header>

        {message && <div className="notice">{message}</div>}

        {active === "dashboard" && (
          <Dashboard
            summary={summary}
            range={analyticsRange}
            onRangeChange={setAnalyticsRange}
            onApply={loadDashboard}
          />
        )}
        {active === "profile" && <ProfilePanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "workout" && <WorkoutPanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "nutrition" && <NutritionPanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "body" && <BodyPanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "goals" && <GoalPanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "reports" && <ReportPanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "ai" && <AIPanel token={token} onChanged={refreshAll} setMessage={setMessage} />}
        {active === "admin" && user?.role === "admin" && (
          <AdminPanel token={token} onChanged={refreshAll} setMessage={setMessage} />
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
    const body = mode === "login" ? { email: form.email, password: form.password } : form;
    const result = await api<AuthResponse>(`/auth/${mode}`, { method: "POST", body: JSON.stringify(body) });
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

function Dashboard({
  summary,
  range,
  onRangeChange,
  onApply
}: {
  summary: DashboardSummary | null;
  range: { period: string; startDate: string; endDate: string };
  onRangeChange: (range: { period: string; startDate: string; endDate: string }) => void;
  onApply: () => void;
}) {
  const cards = useMemo(
    () => [
      ["Workout Sessions", summary?.workoutSessions ?? 0],
      ["Training Volume", `${round(summary?.trainingVolume)} kg`],
      ["Consistency", `${round(summary?.workoutConsistency)}%`],
      ["Progressive Load", readable(summary?.progressiveOverloadStatus)],
      ["Calories In", `${round(summary?.caloriesIn)} kcal`],
      ["Meal Quality", `${round(summary?.mealQualityScore)} / 100`],
      ["Calorie Status", readable(summary?.calorieStatus)],
      ["Weight Trend", readable(summary?.weightTrend) || "No data"],
      ["7-Day Weight Avg", `${round(summary?.weightMovingAverage7DayKg)} kg`],
      ["Goal Adherence", `${round(summary?.goalAdherence)}%`]
    ],
    [summary]
  );

  return (
    <section className="panel-stack">
      <Panel title="Analytics Range">
        <div className="filter-row">
          <label>
            Period
            <select value={range.period} onChange={(event) => onRangeChange({ ...range, period: event.target.value })}>
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
              <option value="custom">Custom range</option>
            </select>
          </label>
          {range.period === "custom" && (
            <>
              <label>
                Start date
                <input
                  type="date"
                  value={range.startDate}
                  onChange={(event) => onRangeChange({ ...range, startDate: event.target.value })}
                />
              </label>
              <label>
                End date
                <input
                  type="date"
                  value={range.endDate}
                  onChange={(event) => onRangeChange({ ...range, endDate: event.target.value })}
                />
              </label>
            </>
          )}
          <button className="primary" onClick={onApply}>Apply</button>
        </div>
        <span className="muted">
          Current range: {summary ? `${formatDate(summary.startDate)} to ${formatDate(summary.endDate)}` : "No data loaded"}
        </span>
      </Panel>
      <div className="content-grid">
        {cards.map(([label, value]) => (
          <article className="metric" key={label}>
            <span>{label}</span>
            <strong>{value}</strong>
          </article>
        ))}
      </div>
      <section className="two-column">
        <Panel title="Macro Split">
          <div className="content-grid compact">
            <article className="metric small">
              <span>Protein Intake</span>
              <strong>{round(summary?.protein)} g</strong>
            </article>
            <article className="metric small">
              <span>Carbs Intake</span>
              <strong>{round(summary?.carbohydrates)} g</strong>
            </article>
            <article className="metric small">
              <span>Fat Intake</span>
              <strong>{round(summary?.fat)} g</strong>
            </article>
          </div>
          <Ratio label="Protein" value={summary?.proteinRatio ?? 0} />
          <Ratio label="Carbs" value={summary?.carbohydrateRatio ?? 0} />
          <Ratio label="Fat" value={summary?.fatRatio ?? 0} />
        </Panel>
        <Panel title="Muscle Groups">
          {Object.entries(summary?.muscleGroupDistribution ?? {}).length === 0 && <EmptyState text="No workout data yet" />}
          {Object.entries(summary?.muscleGroupDistribution ?? {}).map(([group, count]) => (
            <div className="split-row" key={group}>
              <span>{readable(group)}</span>
              <strong>{count}</strong>
            </div>
          ))}
        </Panel>
      </section>
      <section className="two-column">
        <Panel title="Nutrition Gaps">
          <TagList values={summary?.nutritionGaps ?? []} empty="No major gaps detected" />
        </Panel>
        <Panel title="Training Warnings">
          <TagList values={summary?.muscleGroupWarnings ?? []} empty="Training distribution looks balanced" />
        </Panel>
      </section>
    </section>
  );
}

function ProfilePanel({ token, onChanged, setMessage }: PanelProps) {
  const [profile, setProfile] = useState<UserProfile | null>(null);

  useEffect(() => {
    api<UserProfile>("/profile", {}, token).then((result) => {
      if (result.data) setProfile(result.data);
    });
  }, [token]);

  return (
    <Panel title="Profile">
      <SimpleForm
        fields={[
          ["age", "Age", "number", String(profile?.age ?? 22)],
          ["gender", "Gender", "select", profile?.gender ?? "male", options(["male", "female", "other"])],
          ["heightCm", "Height cm", "number", String(profile?.heightCm ?? 175)],
          ["weightKg", "Weight kg", "number", String(profile?.weightKg ?? 70)],
          [
            "activityLevel",
            "Activity level",
            "select",
            profile?.activityLevel ?? "moderate",
            options(["sedentary", "light", "moderate", "active", "very_active"])
          ],
          [
            "primaryGoal",
            "Primary goal",
            "select",
            profile?.primaryGoal ?? "fat_loss",
            options(["fat_loss", "muscle_gain", "maintenance", "strength", "health"])
          ]
        ]}
        submitLabel="Save profile"
        onSubmit={async (values) => {
          const result = await api<UserProfile>("/profile", { method: "PUT", body: JSON.stringify(coerceNumbers(values)) }, token);
          if (result.error) return setMessage(result.error);
          setProfile(result.data ?? null);
          setMessage("Profile saved");
          onChanged();
        }}
      />
    </Panel>
  );
}

function WorkoutPanel({ token, onChanged, setMessage }: PanelProps) {
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [sessions, setSessions] = useState<WorkoutSession[]>([]);
  const [editing, setEditing] = useState<WorkoutSession | null>(null);
  const [exerciseSearch, setExerciseSearch] = useState("");
  const [selectedWorkoutDate, setSelectedWorkoutDate] = useState(today());

  const filteredExercises = useMemo(() => {
    const keyword = exerciseSearch.trim().toLowerCase();
    if (!keyword) return exercises;
    return exercises.filter((exercise) =>
      `${exercise.name} ${exercise.category} ${exercise.muscleGroup} ${exercise.equipment}`.toLowerCase().includes(keyword)
    );
  }, [exercises, exerciseSearch]);

  const sessionsForSelectedDate = useMemo(
    () => sessions.filter((session) => dateInput(session.workoutDate) === selectedWorkoutDate),
    [sessions, selectedWorkoutDate]
  );

  const workoutTotals = useMemo(
    () =>
      sessionsForSelectedDate.reduce(
        (total, session) => ({
          sessions: total.sessions + 1,
          duration: total.duration + session.durationMin,
          volume: total.volume + session.entries.reduce((sum, entry) => sum + workoutEntryVolume(entry), 0),
          entries: total.entries + session.entries.length
        }),
        { sessions: 0, duration: 0, volume: 0, entries: 0 }
      ),
    [sessionsForSelectedDate]
  );

  async function reload() {
    const [exerciseResult, sessionResult] = await Promise.all([
      api<Exercise[]>("/workouts/exercises", {}, token),
      api<WorkoutSession[]>("/workouts", {}, token)
    ]);
    if (exerciseResult.data) setExercises(exerciseResult.data);
    if (sessionResult.data) setSessions(sessionResult.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function remove(id: string) {
    const result = await api(`/workouts/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Workout deleted");
    await reload();
    onChanged();
  }

  return (
    <section className="panel-stack">
      <Panel title="Workout Day">
        <div className="filter-row">
          <label>
            Date
            <input
              type="date"
              value={selectedWorkoutDate}
              onChange={(event) => {
                setSelectedWorkoutDate(event.target.value);
                setEditing(null);
              }}
            />
          </label>
        </div>
        <span className="muted">Workout history and training totals are calculated for the selected date.</span>
      </Panel>
      <div className="content-grid">
        <article className="metric">
          <span>Workout Sessions</span>
          <strong>{workoutTotals.sessions}</strong>
        </article>
        <article className="metric">
          <span>Total Duration</span>
          <strong>{workoutTotals.duration} min</strong>
        </article>
        <article className="metric">
          <span>Training Load</span>
          <strong>{round(workoutTotals.volume)}</strong>
        </article>
        <article className="metric">
          <span>Exercise Entries</span>
          <strong>{workoutTotals.entries}</strong>
        </article>
      </div>
      <Panel title="Log Workout">
        <div className="filter-row">
          <label>
            Search exercise
            <input
              value={exerciseSearch}
              onChange={(event) => setExerciseSearch(event.target.value)}
              placeholder="Search by exercise, muscle group, or equipment"
            />
          </label>
        </div>
        <SimpleForm
          fields={[
            ["workoutDate", "Workout date", "date", dateInput(editing?.workoutDate) || selectedWorkoutDate],
            ["category", "Category", "select", editing?.category ?? "strength", options(["strength", "cardio", "mobility"])],
            ["durationMin", "Duration min", "number", String(editing?.durationMin ?? 45)],
            [
              "exerciseId",
              "Exercise",
              "select",
              editing?.entries?.[0]?.exerciseId ?? filteredExercises[0]?.id ?? "",
              filteredExercises.map((exercise) => ({
                value: exercise.id,
                label: `${exercise.name} - ${readable(exercise.muscleGroup)}`
              }))
            ],
            ["sets", "Sets", "number", String(editing?.entries?.[0]?.sets ?? 3)],
            ["reps", "Reps", "number", String(editing?.entries?.[0]?.reps ?? 10)],
            ["weightKg", "Weight kg", "number", String(editing?.entries?.[0]?.weightKg ?? 40)],
            ["restSec", "Rest sec", "number", String(editing?.entries?.[0]?.restSec ?? 60)],
            ["notes", "Notes", "text", editing?.notes ?? ""]
          ]}
          submitLabel={editing ? "Update workout" : "Add workout"}
          onSubmit={async (values) => {
            const payload = {
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
            };
            const result = await api(editing ? `/workouts/${editing.id}` : "/workouts", {
              method: editing ? "PUT" : "POST",
              body: JSON.stringify(payload)
            }, token);
            if (result.error) return setMessage(result.error);
            setMessage(editing ? "Workout updated" : "Workout logged");
            setEditing(null);
            await reload();
            onChanged();
          }}
        />
        {editing && <button onClick={() => setEditing(null)}>Cancel edit</button>}
      </Panel>
      <Panel title={`Workout History - ${selectedWorkoutDate}`}>
        <DataTable
          headers={["Date", "Category", "Duration", "Exercises", "Actions"]}
          empty="No workouts for selected date"
          rows={sessionsForSelectedDate.map((session) => [
            formatDate(session.workoutDate),
            readable(session.category),
            `${session.durationMin} min`,
            session.entries.map(formatWorkoutEntry).join(", "),
            <div className="action-row" key={session.id}>
              <button onClick={() => setEditing(session)}>Edit</button>
              <IconAction title="Delete workout" onClick={() => remove(session.id)} icon={<Trash2 size={16} />} />
            </div>
          ])}
        />
      </Panel>
      <Panel title="Exercise Reference">
        <DataTable
          headers={["Exercise", "Category", "Muscle Group", "Equipment", "Intensity"]}
          empty="No exercises match the search"
          rows={filteredExercises.map((exercise) => [
            exercise.name,
            readable(exercise.category),
            readable(exercise.muscleGroup),
            readable(exercise.equipment),
            readable(exercise.intensityLevel)
          ])}
        />
      </Panel>
    </section>
  );
}

function NutritionPanel({ token, onChanged, setMessage }: PanelProps) {
  const [foods, setFoods] = useState<FoodItem[]>([]);
  const [meals, setMeals] = useState<MealLog[]>([]);
  const [editing, setEditing] = useState<MealLog | null>(null);
  const [foodSearch, setFoodSearch] = useState("");
  const [selectedMealDate, setSelectedMealDate] = useState(today());

  const filteredFoods = useMemo(() => {
    const keyword = foodSearch.trim().toLowerCase();
    if (!keyword) return foods;
    return foods.filter((food) => `${food.name} ${food.servingSize}`.toLowerCase().includes(keyword));
  }, [foods, foodSearch]);

  const mealsForSelectedDate = useMemo(
    () => meals.filter((meal) => dateInput(meal.mealTime) === selectedMealDate),
    [meals, selectedMealDate]
  );

  const usedFoods = useMemo(() => {
    const map = new Map<string, FoodItem>();
    for (const meal of mealsForSelectedDate) {
      if (meal.foodItem) map.set(meal.foodItem.id, meal.foodItem);
    }
    return Array.from(map.values());
  }, [mealsForSelectedDate]);

  const nutritionTotals = useMemo(
    () =>
      mealsForSelectedDate.reduce(
        (total, meal) => ({
          calories: total.calories + meal.totalCalories,
          protein: total.protein + meal.totalProtein,
          carbs: total.carbs + meal.totalCarbs,
          fat: total.fat + meal.totalFat
        }),
        { calories: 0, protein: 0, carbs: 0, fat: 0 }
      ),
    [mealsForSelectedDate]
  );

  async function reload() {
    const [foodResult, mealResult] = await Promise.all([
      api<FoodItem[]>("/nutrition/foods", {}, token),
      api<MealLog[]>("/nutrition/meals", {}, token)
    ]);
    if (foodResult.data) setFoods(foodResult.data);
    if (mealResult.data) setMeals(mealResult.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function remove(id: string) {
    const result = await api(`/nutrition/meals/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Meal deleted");
    await reload();
    onChanged();
  }

  return (
    <section className="panel-stack">
      <Panel title="Nutrition Day">
        <div className="filter-row">
          <label>
            Date
            <input
              type="date"
              value={selectedMealDate}
              onChange={(event) => {
                setSelectedMealDate(event.target.value);
                setEditing(null);
              }}
            />
          </label>
        </div>
        <span className="muted">Meal history and nutrition totals are calculated for the selected date.</span>
      </Panel>
      <div className="content-grid">
        <article className="metric">
          <span>Total Calories</span>
          <strong>{round(nutritionTotals.calories)} kcal</strong>
        </article>
        <article className="metric">
          <span>Total Protein</span>
          <strong>{round(nutritionTotals.protein)} g</strong>
        </article>
        <article className="metric">
          <span>Total Carbs</span>
          <strong>{round(nutritionTotals.carbs)} g</strong>
        </article>
        <article className="metric">
          <span>Total Fat</span>
          <strong>{round(nutritionTotals.fat)} g</strong>
        </article>
      </div>
      <Panel title="Log Meal">
        <div className="filter-row">
          <label>
            Search food
            <input
              value={foodSearch}
              onChange={(event) => setFoodSearch(event.target.value)}
              placeholder="Search by food name or serving size"
            />
          </label>
        </div>
        <SimpleForm
          fields={[
            [
              "foodItemId",
              "Food",
              "select",
              editing?.foodItemId ?? filteredFoods[0]?.id ?? "",
              filteredFoods.map((food) => ({
                value: food.id,
                label: `${food.name} - ${food.servingSize}`
              }))
            ],
            ["mealType", "Meal type", "select", editing?.mealType ?? "lunch", options(["breakfast", "lunch", "dinner", "snack"])],
            ["quantity", "Quantity", "number", String(editing?.quantity ?? 1)],
            ["mealTime", "Meal time", "datetime-local", dateTimeInput(editing?.mealTime) || localDateTimeForDate(selectedMealDate)]
          ]}
          submitLabel={editing ? "Update meal" : "Add meal"}
          onSubmit={async (values) => {
            const payload = {
              foodItemId: values.foodItemId,
              mealType: values.mealType,
              quantity: Number(values.quantity),
              mealTime: new Date(values.mealTime).toISOString()
            };
            const result = await api(editing ? `/nutrition/meals/${editing.id}` : "/nutrition/meals", {
              method: editing ? "PUT" : "POST",
              body: JSON.stringify(payload)
            }, token);
            if (result.error) return setMessage(result.error);
            setMessage(editing ? "Meal updated" : "Meal logged");
            setEditing(null);
            await reload();
            onChanged();
          }}
        />
        {editing && <button onClick={() => setEditing(null)}>Cancel edit</button>}
      </Panel>
      <Panel title={`Meal History - ${selectedMealDate}`}>
        <DataTable
          headers={["Time", "Meal", "Food", "Calories", "Protein", "Actions"]}
          empty="No meals for selected date"
          rows={mealsForSelectedDate.map((meal) => [
            formatDateTime(meal.mealTime),
            readable(meal.mealType),
            meal.foodItem?.name ?? "Food",
            `${round(meal.totalCalories)} kcal`,
            `${round(meal.totalProtein)} g`,
            <div className="action-row" key={meal.id}>
              <button onClick={() => setEditing(meal)}>Edit</button>
              <IconAction title="Delete meal" onClick={() => remove(meal.id)} icon={<Trash2 size={16} />} />
            </div>
          ])}
        />
      </Panel>
      <Panel title="Used Food Reference">
        <DataTable
          headers={["Food", "Serving", "Calories", "Protein", "Carbs", "Fat"]}
          empty="No food has been used in meal history yet"
          rows={usedFoods.map((food) => [
            food.name,
            food.servingSize,
            `${round(food.calories)} kcal`,
            `${round(food.protein)} g`,
            `${round(food.carbohydrates)} g`,
            `${round(food.fat)} g`
          ])}
        />
      </Panel>
    </section>
  );
}

function BodyPanel({ token, onChanged, setMessage }: PanelProps) {
  const [records, setRecords] = useState<BodyRecord[]>([]);
  const [editing, setEditing] = useState<BodyRecord | null>(null);

  async function reload() {
    const result = await api<BodyRecord[]>("/body-records", {}, token);
    if (result.data) setRecords(result.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function remove(id: string) {
    const result = await api(`/body-records/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Body record deleted");
    await reload();
    onChanged();
  }

  return (
    <section className="panel-stack">
      <Panel title="Add Body Record">
        <SimpleForm
          fields={[
            ["recordDate", "Record date", "date", dateInput(editing?.recordDate) || today()],
            ["weightKg", "Weight kg", "number", String(editing?.weightKg ?? 70)],
            ["note", "Note", "text", editing?.note ?? ""]
          ]}
          submitLabel={editing ? "Update record" : "Add record"}
          onSubmit={async (values) => {
            const result = await api(editing ? `/body-records/${editing.id}` : "/body-records", {
              method: editing ? "PUT" : "POST",
              body: JSON.stringify({ recordDate: values.recordDate, weightKg: Number(values.weightKg), note: values.note })
            }, token);
            if (result.error) return setMessage(result.error);
            setMessage(editing ? "Body record updated" : "Body record added");
            setEditing(null);
            await reload();
            onChanged();
          }}
        />
        {editing && <button onClick={() => setEditing(null)}>Cancel edit</button>}
      </Panel>
      <Panel title="Body Timeline">
        <DataTable
          headers={["Date", "Weight", "Note", "Actions"]}
          empty="No body records yet"
          rows={records.map((record) => [
            formatDate(record.recordDate),
            `${round(record.weightKg)} kg`,
            record.note || "-",
            <div className="action-row" key={record.id}>
              <button onClick={() => setEditing(record)}>Edit</button>
              <IconAction title="Delete body record" onClick={() => remove(record.id)} icon={<Trash2 size={16} />} />
            </div>
          ])}
        />
      </Panel>
    </section>
  );
}

function GoalPanel({ token, onChanged, setMessage }: PanelProps) {
  const [goals, setGoals] = useState<Goal[]>([]);
  const [editing, setEditing] = useState<Goal | null>(null);

  async function reload() {
    const result = await api<Goal[]>("/goals", {}, token);
    if (result.data) setGoals(result.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function remove(id: string) {
    const result = await api(`/goals/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Goal deleted");
    await reload();
    onChanged();
  }

  async function updateStatus(goal: Goal, status: string) {
    const result = await api<Goal>(`/goals/${goal.id}/status`, { method: "PATCH", body: JSON.stringify({ status }) }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Goal status updated");
    await reload();
    onChanged();
  }

  async function updateMilestone(goal: Goal, milestone: GoalMilestone) {
    const result = await api(`/goals/${goal.id}/milestones/${milestone.id}`, {
      method: "PATCH",
      body: JSON.stringify({ completed: !milestone.completed })
    }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Milestone updated");
    await reload();
    onChanged();
  }

  return (
    <section className="panel-stack">
      <Panel title="Create Goal">
        <SimpleForm
          fields={[
            ["goalType", "Goal type", "select", editing?.goalType ?? "workout_frequency", options(["workout_frequency", "body_weight", "nutrition", "strength"])],
            ["targetMetric", "Target metric", "text", editing?.targetMetric ?? "sessions_per_week"],
            ["targetValue", "Target value", "number", String(editing?.targetValue ?? 3)],
            ["deadline", "Deadline", "date", dateInput(editing?.deadline) || nextMonth()],
            ["priority", "Priority", "select", editing?.priority ?? "high", options(["low", "medium", "high"])],
            ["milestone1Title", "Milestone 1 title", "text", editing?.milestones?.[0]?.title ?? "Milestone 1"],
            ["milestone1Target", "Milestone 1 target", "number", String(editing?.milestones?.[0]?.targetValue ?? 1)],
            ["milestone1DueDate", "Milestone 1 due date", "date", dateInput(editing?.milestones?.[0]?.dueDate) || nextWeek(1)],
            ["milestone2Title", "Milestone 2 title", "text", editing?.milestones?.[1]?.title ?? "Milestone 2"],
            ["milestone2Target", "Milestone 2 target", "number", String(editing?.milestones?.[1]?.targetValue ?? 2)],
            ["milestone2DueDate", "Milestone 2 due date", "date", dateInput(editing?.milestones?.[1]?.dueDate) || nextWeek(2)],
            ["milestone3Title", "Milestone 3 title", "text", editing?.milestones?.[2]?.title ?? "Milestone 3"],
            ["milestone3Target", "Milestone 3 target", "number", String(editing?.milestones?.[2]?.targetValue ?? 3)],
            ["milestone3DueDate", "Milestone 3 due date", "date", dateInput(editing?.milestones?.[2]?.dueDate) || nextWeek(4)]
          ]}
          submitLabel={editing ? "Update goal" : "Create goal"}
          onSubmit={async (values) => {
            const milestones = [1, 2, 3].map((index) => ({
              title: values[`milestone${index}Title`],
              targetValue: Number(values[`milestone${index}Target`]),
              dueDate: values[`milestone${index}DueDate`]
            }));
            const result = await api(editing ? `/goals/${editing.id}` : "/goals", {
              method: editing ? "PUT" : "POST",
              body: JSON.stringify({
                goalType: values.goalType,
                targetMetric: values.targetMetric,
                targetValue: Number(values.targetValue),
                deadline: values.deadline,
                priority: values.priority,
                milestones
              })
            }, token);
            if (result.error) return setMessage(result.error);
            setMessage(editing ? "Goal updated" : "Goal created");
            setEditing(null);
            await reload();
            onChanged();
          }}
        />
        {editing && <button onClick={() => setEditing(null)}>Cancel edit</button>}
      </Panel>
      {goals.length === 0 && <EmptyState text="No goals yet" />}
      {goals.map((goal) => (
        <Panel key={goal.id} title={`${readable(goal.goalType)} - ${readable(goal.status)}`}>
          <div className="goal-card">
            <div>
              <span className="muted">{goal.targetMetric}</span>
              <strong>{goal.targetValue}</strong>
              <span className="muted">Deadline {formatDate(goal.deadline)} / {readable(goal.priority)}</span>
            </div>
            <div className="action-row">
              <button onClick={() => setEditing(goal)}>Edit</button>
              <button onClick={() => updateStatus(goal, goal.status === "active" ? "completed" : "active")}>
                {goal.status === "active" ? "Complete" : "Reactivate"}
              </button>
              <IconAction title="Delete goal" onClick={() => remove(goal.id)} icon={<Trash2 size={16} />} />
            </div>
          </div>
          {(goal.milestones ?? []).length > 0 && (
            <div className="milestone-list">
              {(goal.milestones ?? []).map((milestone) => (
                <button
                  className={milestone.completed ? "milestone complete" : "milestone"}
                  key={milestone.id}
                  onClick={() => updateMilestone(goal, milestone)}
                >
                  <CheckCircle size={16} />
                  <span>{milestone.title}</span>
                </button>
              ))}
            </div>
          )}
        </Panel>
      ))}
    </section>
  );
}

function ReportPanel({ token, onChanged, setMessage }: PanelProps) {
  const [reports, setReports] = useState<ProgressReport[]>([]);
  const [range, setRange] = useState({ periodType: "weekly", startDate: today(), endDate: today() });

  async function reload() {
    const result = await api<ProgressReport[]>("/reports", {}, token);
    if (result.data) setReports(result.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function generate() {
    const payload =
      range.periodType === "custom"
        ? range
        : { periodType: range.periodType };
    const result = await api<{ report: ProgressReport; metrics: DashboardSummary }>(
      "/reports/generate",
      { method: "POST", body: JSON.stringify(payload) },
      token
    );
    if (result.error) return setMessage(result.error);
    setMessage(`${range.periodType} report generated`);
    await reload();
    onChanged();
  }

  async function remove(id: string) {
    const result = await api(`/reports/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Report deleted");
    await reload();
  }

  async function download(report: ProgressReport) {
    const result = await downloadFile(`/reports/${report.id}/download`, token);
    if (result.error || !result.data) return setMessage(result.error ?? "Download failed");
    const url = URL.createObjectURL(result.data);
    const link = document.createElement("a");
    link.href = url;
    link.download = `progress-report-${report.id}.html`;
    link.click();
    URL.revokeObjectURL(url);
  }

  return (
    <section className="panel-stack">
      <Panel title="Generate Report">
        <div className="filter-row">
          <label>
            Period
            <select value={range.periodType} onChange={(event) => setRange({ ...range, periodType: event.target.value })}>
              <option value="daily">Daily</option>
              <option value="weekly">Weekly</option>
              <option value="monthly">Monthly</option>
              <option value="custom">Custom range</option>
            </select>
          </label>
          {range.periodType === "custom" && (
            <>
              <label>
                Start date
                <input type="date" value={range.startDate} onChange={(event) => setRange({ ...range, startDate: event.target.value })} />
              </label>
              <label>
                End date
                <input type="date" value={range.endDate} onChange={(event) => setRange({ ...range, endDate: event.target.value })} />
              </label>
            </>
          )}
          <button className="primary" onClick={generate}>Generate report</button>
        </div>
      </Panel>
      {reports.length === 0 && <EmptyState text="No reports yet" />}
      {reports.map((report) => (
        <article className="recommendation" key={report.id}>
          <span>{readable(report.periodType)} / {formatDateTime(report.generatedAt)}</span>
          <pre>{report.summary}</pre>
          <div className="action-row">
            <IconAction title="Download report" onClick={() => download(report)} icon={<Download size={16} />} label="Download" />
            <IconAction title="Delete report" onClick={() => remove(report.id)} icon={<Trash2 size={16} />} label="Delete" />
          </div>
        </article>
      ))}
    </section>
  );
}

function AIPanel({ token, onChanged, setMessage }: PanelProps) {
  const [recommendations, setRecommendations] = useState<Recommendation[]>([]);

  async function reload() {
    const result = await api<Recommendation[]>("/recommendations", {}, token);
    if (result.data) setRecommendations(result.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function generate(type: string) {
    const result = await api<Recommendation>("/recommendations/generate", { method: "POST", body: JSON.stringify({ type }) }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Recommendation generated");
    await reload();
    onChanged();
  }

  async function remove(id: string) {
    const result = await api(`/recommendations/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Recommendation deleted");
    await reload();
  }

  async function feedback(id: string, rating: string) {
    const result = await api(`/recommendations/${id}/feedback`, {
      method: "POST",
      body: JSON.stringify({ rating, suitability: "general", comment: "" })
    }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Feedback submitted");
  }

  return (
    <section className="panel-stack">
      <div className="action-row">
        <button className="primary" onClick={() => generate("weekly")}>Generate weekly feedback</button>
        <button onClick={() => generate("workout")}>Workout plan</button>
        <button onClick={() => generate("meal")}>Meal idea</button>
      </div>
      {recommendations.length === 0 && <EmptyState text="No AI recommendations yet" />}
      {recommendations.map((item) => (
        <article className="recommendation" key={item.id}>
          <span>{readable(item.type)} / {formatDateTime(item.createdAt)}</span>
          <pre>{item.content}</pre>
          <div className="action-row">
            <button onClick={() => feedback(item.id, "useful")}>Useful</button>
            <button onClick={() => feedback(item.id, "not_useful")}>Not useful</button>
            <IconAction title="Delete recommendation" onClick={() => remove(item.id)} icon={<Trash2 size={16} />} label="Delete" />
          </div>
        </article>
      ))}
    </section>
  );
}

function AdminPanel({ token, setMessage }: PanelProps) {
  const [users, setUsers] = useState<User[]>([]);
  const [foods, setFoods] = useState<FoodItem[]>([]);
  const [exercises, setExercises] = useState<Exercise[]>([]);
  const [editingFood, setEditingFood] = useState<FoodItem | null>(null);
  const [editingExercise, setEditingExercise] = useState<Exercise | null>(null);

  async function reload() {
    const [usersResult, foodsResult, exerciseResult] = await Promise.all([
      api<User[]>("/admin/users", {}, token),
      api<FoodItem[]>("/nutrition/foods", {}, token),
      api<Exercise[]>("/workouts/exercises", {}, token)
    ]);
    if (usersResult.data) setUsers(usersResult.data);
    if (foodsResult.data) setFoods(foodsResult.data);
    if (exerciseResult.data) setExercises(exerciseResult.data);
  }

  useEffect(() => {
    reload();
  }, [token]);

  async function updateRole(user: User, role: string) {
    const result = await api<User>(`/admin/users/${user.id}/role`, { method: "PATCH", body: JSON.stringify({ role }) }, token);
    if (result.error) return setMessage(result.error);
    setMessage("User role updated");
    await reload();
  }

  async function deleteFood(id: string) {
    const result = await api(`/admin/nutrition/foods/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Food deleted");
    await reload();
  }

  async function deleteExercise(id: string) {
    const result = await api(`/admin/workouts/exercises/${id}`, { method: "DELETE" }, token);
    if (result.error) return setMessage(result.error);
    setMessage("Exercise deleted");
    await reload();
  }

  return (
    <section className="panel-stack">
      <Panel title="Users">
        <DataTable
          headers={["Name", "Email", "Role", "Actions"]}
          empty="No users"
          rows={users.map((item) => [
            item.name,
            item.email,
            <StatusPill key={`${item.id}-role`} value={item.role} />,
            <button key={item.id} onClick={() => updateRole(item, item.role === "admin" ? "user" : "admin")}>
              Make {item.role === "admin" ? "user" : "admin"}
            </button>
          ])}
        />
      </Panel>
      <section className="two-column">
        <Panel title="Add Food">
          <SimpleForm
            fields={[
              ["name", "Name", "text", editingFood?.name ?? "Chicken breast"],
              ["servingSize", "Serving size", "text", editingFood?.servingSize ?? "100g"],
              ["calories", "Calories", "number", String(editingFood?.calories ?? 165)],
              ["protein", "Protein", "number", String(editingFood?.protein ?? 31)],
              ["carbohydrates", "Carbs", "number", String(editingFood?.carbohydrates ?? 0)],
              ["fat", "Fat", "number", String(editingFood?.fat ?? 3.6)],
              ["sugar", "Sugar", "number", String(editingFood?.sugar ?? 0)],
              ["sodium", "Sodium", "number", String(editingFood?.sodium ?? 74)]
            ]}
            submitLabel={editingFood ? "Update food" : "Add food"}
            onSubmit={async (values) => {
              const result = await api(editingFood ? `/admin/nutrition/foods/${editingFood.id}` : "/admin/nutrition/foods", {
                method: editingFood ? "PUT" : "POST",
                body: JSON.stringify(coerceNumbers(values))
              }, token);
              if (result.error) return setMessage(result.error);
              setMessage(editingFood ? "Food updated" : "Food added");
              setEditingFood(null);
              await reload();
            }}
          />
          {editingFood && <button onClick={() => setEditingFood(null)}>Cancel edit</button>}
        </Panel>
        <Panel title="Add Exercise">
          <SimpleForm
            fields={[
              ["name", "Name", "text", editingExercise?.name ?? "Overhead Press"],
              ["category", "Category", "select", editingExercise?.category ?? "strength", options(["strength", "cardio", "mobility"])],
              ["muscleGroup", "Muscle group", "text", editingExercise?.muscleGroup ?? "shoulders"],
              ["equipment", "Equipment", "text", editingExercise?.equipment ?? "barbell"],
              ["intensityLevel", "Intensity", "select", editingExercise?.intensityLevel ?? "moderate", options(["low", "moderate", "high", "variable"])]
            ]}
            submitLabel={editingExercise ? "Update exercise" : "Add exercise"}
            onSubmit={async (values) => {
              const result = await api(editingExercise ? `/admin/workouts/exercises/${editingExercise.id}` : "/admin/workouts/exercises", {
                method: editingExercise ? "PUT" : "POST",
                body: JSON.stringify(values)
              }, token);
              if (result.error) return setMessage(result.error);
              setMessage(editingExercise ? "Exercise updated" : "Exercise added");
              setEditingExercise(null);
              await reload();
            }}
          />
          {editingExercise && <button onClick={() => setEditingExercise(null)}>Cancel edit</button>}
        </Panel>
      </section>
      <section className="two-column">
        <Panel title="Food Library">
          <DataTable
            headers={["Food", "Calories", "Actions"]}
            empty="No foods"
            rows={foods.map((food) => [
              food.name,
              `${round(food.calories)} kcal`,
              <div className="action-row" key={food.id}>
                <button onClick={() => setEditingFood(food)}>Edit</button>
                <IconAction title="Delete food" onClick={() => deleteFood(food.id)} icon={<Trash2 size={16} />} />
              </div>
            ])}
          />
        </Panel>
        <Panel title="Exercise Library">
          <DataTable
            headers={["Exercise", "Group", "Actions"]}
            empty="No exercises"
            rows={exercises.map((exercise) => [
              exercise.name,
              readable(exercise.muscleGroup),
              <div className="action-row" key={exercise.id}>
                <button onClick={() => setEditingExercise(exercise)}>Edit</button>
                <IconAction title="Delete exercise" onClick={() => deleteExercise(exercise.id)} icon={<Trash2 size={16} />} />
              </div>
            ])}
          />
        </Panel>
      </section>
    </section>
  );
}

type Field = [key: string, label: string, type: string, value: string, options?: SelectOption[]];

type SelectOption = {
  value: string;
  label: string;
};

function SimpleForm({
  fields,
  submitLabel,
  onSubmit
}: {
  fields: Field[];
  submitLabel: string;
  onSubmit: (values: Record<string, string>) => Promise<void>;
}) {
  const initial = Object.fromEntries(fields.map(([key, , , value]) => [key, value]));
  const [values, setValues] = useState<Record<string, string>>(initial);

  useEffect(() => {
    setValues(Object.fromEntries(fields.map(([key, , , value]) => [key, value])));
  }, [fields.map((field) => `${field[0]}:${field[3]}:${field[4]?.length ?? 0}`).join("|")]);

  async function submit(event: FormEvent) {
    event.preventDefault();
    await onSubmit(values);
  }

  return (
    <form className="data-form" onSubmit={submit}>
      {fields.map(([key, label, type, , fieldOptions]) => (
        <label key={key}>
          {label}
          {type === "select" ? (
            <select value={values[key] ?? ""} onChange={(event) => setValues({ ...values, [key]: event.target.value })} required>
              <option value="">Select</option>
              {(fieldOptions ?? []).map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          ) : (
            <input
              type={type}
              value={values[key] ?? ""}
              onChange={(event) => setValues({ ...values, [key]: event.target.value })}
              required={!["notes", "note", "comment"].includes(key)}
            />
          )}
        </label>
      ))}
      <button className="primary" type="submit">{submitLabel}</button>
    </form>
  );
}

function Panel({ title, children }: { title: string; children: ReactNode }) {
  return (
    <section className="panel">
      <h2>{title}</h2>
      {children}
    </section>
  );
}

function DataTable({ headers, rows, empty }: { headers: string[]; rows: ReactNode[][]; empty: string }) {
  if (rows.length === 0) return <EmptyState text={empty} />;
  return (
    <div className="table-wrap">
      <table>
        <thead>
          <tr>{headers.map((header) => <th key={header}>{header}</th>)}</tr>
        </thead>
        <tbody>
          {rows.map((row, index) => (
            <tr key={index}>
              {row.map((cell, cellIndex) => <td key={cellIndex}>{cell}</td>)}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

function Ratio({ label, value }: { label: string; value: number }) {
  return (
    <div className="ratio">
      <div className="split-row">
        <span>{label}</span>
        <strong>{round(value)}%</strong>
      </div>
      <div className="bar">
        <div style={{ width: `${Math.min(Math.max(value, 0), 100)}%` }} />
      </div>
    </div>
  );
}

function TagList({ values, empty }: { values: string[]; empty: string }) {
  if (values.length === 0) return <EmptyState text={empty} />;
  return (
    <div className="tag-list">
      {values.map((value) => <StatusPill key={value} value={readable(value)} />)}
    </div>
  );
}

function StatusPill({ value }: { value: string }) {
  return <span className={`pill ${value.toLowerCase().replace(/ /g, "-")}`}>{readable(value)}</span>;
}

function formatWorkoutEntry(entry: WorkoutEntry) {
  const weight = entry.weightKg > 0 ? ` @ ${round(entry.weightKg)}kg` : "";
  return `${entry.exercise?.name ?? "Exercise"} ${entry.sets}x${entry.reps}${weight}`;
}

function workoutEntryVolume(entry: WorkoutEntry) {
  const load = entry.weightKg > 0 ? entry.weightKg : 1;
  return entry.sets * entry.reps * load;
}

function EmptyState({ text }: { text: string }) {
  return <div className="empty-state">{text}</div>;
}

function IconAction({ title, onClick, icon, label }: { title: string; onClick: () => void; icon: ReactNode; label?: string }) {
  return (
    <button className={label ? "action-button" : "icon-button"} onClick={onClick} title={title} type="button">
      {icon}
      {label && <span>{label}</span>}
    </button>
  );
}

function options(values: string[]): SelectOption[] {
  return values.map((value) => ({ value, label: readable(value) }));
}

function dashboardQuery(range: { period: string; startDate: string; endDate: string }) {
  const params = new URLSearchParams({ period: range.period });
  if (range.period === "custom") {
    params.set("startDate", range.startDate);
    params.set("endDate", range.endDate);
  }
  return `?${params.toString()}`;
}

function coerceNumbers(values: Record<string, string>) {
  return Object.fromEntries(
    Object.entries(values).map(([key, value]) => {
      if (value !== "" && !Number.isNaN(Number(value))) return [key, Number(value)];
      return [key, value];
    })
  );
}

function today() {
  return new Date().toISOString().slice(0, 10);
}

function nextMonth() {
  const date = new Date();
  date.setMonth(date.getMonth() + 1);
  return date.toISOString().slice(0, 10);
}

function nextWeek(weeks: number) {
  const date = new Date();
  date.setDate(date.getDate() + weeks * 7);
  return date.toISOString().slice(0, 10);
}

function localDateTime() {
  const date = new Date();
  date.setMinutes(date.getMinutes() - date.getTimezoneOffset());
  return date.toISOString().slice(0, 16);
}

function localDateTimeForDate(dateValue: string) {
  const now = new Date();
  const time = `${String(now.getHours()).padStart(2, "0")}:${String(now.getMinutes()).padStart(2, "0")}`;
  return `${dateValue}T${time}`;
}

function formatDate(value?: string) {
  if (!value) return "-";
  return new Date(value).toLocaleDateString();
}

function formatDateTime(value?: string) {
  if (!value) return "-";
  return new Date(value).toLocaleString();
}

function dateInput(value?: string) {
  if (!value) return "";
  return new Date(value).toISOString().slice(0, 10);
}

function dateTimeInput(value?: string) {
  if (!value) return "";
  const date = new Date(value);
  date.setMinutes(date.getMinutes() - date.getTimezoneOffset());
  return date.toISOString().slice(0, 16);
}

function readable(value?: string) {
  return (value ?? "").replace(/_/g, " ").replace(/\b\w/g, (char: string) => char.toUpperCase());
}

function round(value?: number) {
  return Math.round(value ?? 0);
}
