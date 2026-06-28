package resume

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
)

// TechSkillDictionary contains 200+ technology keywords organized by category
var TechSkillDictionary = map[string][]string{
	"languages": {
		"golang", "go", "java", "python", "javascript", "typescript", "kotlin",
		"swift", "rust", "c++", "c#", "ruby", "php", "scala", "r", "dart",
		"elixir", "clojure", "haskell", "lua", "perl", "groovy", "julia",
		"assembly", "cobol", "fortran", "erlang", "nim",
	},
	"frontend": {
		"react", "reactjs", "react.js", "vue", "vuejs", "vue.js", "angular",
		"angularjs", "svelte", "next.js", "nextjs", "nuxt.js", "gatsby",
		"html", "html5", "css", "css3", "sass", "scss", "less", "tailwindcss",
		"tailwind", "bootstrap", "material-ui", "mui", "chakra ui", "ant design",
		"jquery", "redux", "mobx", "zustand", "recoil", "webpack", "vite",
		"babel", "eslint", "prettier", "storybook", "framer motion",
	},
	"backend": {
		"node.js", "nodejs", "express", "expressjs", "fastapi", "django",
		"flask", "spring", "spring boot", "springboot", "gin", "echo", "fiber",
		"rails", "laravel", "symfony", "nestjs", "nest.js", "fastify",
		"graphql", "rest", "grpc", "soap", "microservices", "api gateway",
		"asp.net", "dotnet", ".net", "quarkus", "micronaut",
	},
	"databases": {
		"postgresql", "postgres", "mysql", "mongodb", "mongodb", "redis",
		"elasticsearch", "cassandra", "dynamodb", "sqlite", "oracle",
		"mssql", "sql server", "mariadb", "cockroachdb", "clickhouse",
		"neo4j", "influxdb", "timescaledb", "firebase", "supabase",
		"planetscale", "vitess", "couchdb", "hbase",
	},
	"cloud": {
		"aws", "amazon web services", "gcp", "google cloud", "azure",
		"ec2", "s3", "lambda", "rds", "eks", "ecs", "cloudformation",
		"terraform", "pulumi", "cdk", "cloud run", "gke", "aks",
		"azure devops", "azure functions", "cloudfront", "route53",
		"sns", "sqs", "kinesis", "eventbridge", "api gateway",
	},
	"devops": {
		"docker", "kubernetes", "k8s", "helm", "jenkins", "github actions",
		"gitlab ci", "circleci", "travis ci", "ansible", "chef", "puppet",
		"prometheus", "grafana", "datadog", "newrelic", "splunk",
		"elk stack", "logstash", "kibana", "nginx", "apache", "istio",
		"envoy", "linkerd", "service mesh", "argocd", "flux",
	},
	"messaging": {
		"kafka", "apache kafka", "rabbitmq", "activemq", "nats", "pulsar",
		"sqs", "pubsub", "eventbus", "zeromq", "zmq",
	},
	"data": {
		"spark", "apache spark", "hadoop", "hive", "flink", "airflow",
		"dbt", "snowflake", "redshift", "bigquery", "databricks",
		"pandas", "numpy", "scikit-learn", "tensorflow", "pytorch",
		"jupyter", "tableau", "power bi", "looker", "dbt",
	},
	"testing": {
		"junit", "testng", "pytest", "jest", "mocha", "chai", "cypress",
		"selenium", "playwright", "postman", "jmeter", "k6",
		"testing library", "vitest", "go test", "mockito", "testify",
	},
	"architecture": {
		"microservices", "monolith", "serverless", "event-driven", "cqrs",
		"event sourcing", "domain driven design", "ddd", "clean architecture",
		"hexagonal architecture", "solid", "design patterns", "mvc",
		"mvvm", "saga pattern", "circuit breaker", "cdc",
	},
	"tools": {
		"git", "github", "gitlab", "bitbucket", "jira", "confluence",
		"slack", "vscode", "intellij", "vim", "linux", "unix",
		"bash", "shell", "powershell", "makefile", "ci/cd",
	},
}

// Certification patterns
var certificationPatterns = []string{
	`(?i)(aws|gcp|azure|google cloud)\s*(certified|certification)`,
	`(?i)(certified|certification)\s*(kubernetes|k8s|docker|devops|developer|architect|sysops|solutions architect|cloud practitioner|developer associate|saa-c|sap-c)`,
	`(?i)cka\b`, `(?i)ckad\b`, `(?i)cks\b`,
	`(?i)pmp\b`, `(?i)scrum master`, `(?i)csm\b`, `(?i)psm\b`,
	`(?i)oracle certified`, `(?i)red hat certified`,
	`(?i)cisco certified`, `(?i)ccna\b`, `(?i)ccnp\b`,
}

// Company patterns (known tech companies)
var knownCompanies = []string{
	"google", "microsoft", "amazon", "apple", "meta", "facebook", "netflix",
	"uber", "airbnb", "stripe", "atlassian", "salesforce", "oracle",
	"ibm", "intel", "nvidia", "adobe", "linkedin", "twitter", "spotify",
	"shopify", "twilio", "okta", "datadog", "snowflake", "confluent",
	"hashicorp", "elastic", "mongodb inc", "redis labs", "cloudera",
	"infosys", "tcs", "wipro", "hcl", "cognizant", "accenture",
	"capgemini", "mindtree", "mphasis", "l&t technology", "tech mahindra",
	"flipkart", "zomato", "swiggy", "paytm", "phonepe", "razorpay",
	"freshworks", "zoho", "byju's", "unacademy", "cred", "meesho",
}

// Education keywords
var educationKeywords = []string{
	"b.tech", "btech", "b.e", "be ", "bachelor", "b.sc", "bsc",
	"m.tech", "mtech", "m.e", "master", "m.sc", "msc", "mba",
	"phd", "ph.d", "doctorate", "b.ca", "mca",
	"university", "college", "institute", "iit", "nit", "iiit",
	"bits pilani", "vit", "manipal", "anna university",
}

// ParseResult holds all extracted resume data
type ParseResult struct {
	Skills         []string                 `json:"skills"`
	Companies      []string                 `json:"companies"`
	Projects       []string                 `json:"projects"`
	Education      []map[string]interface{} `json:"education"`
	Experience     []map[string]interface{} `json:"experience"`
	Certifications []string                 `json:"certifications"`
}

// Parser is the deterministic resume parser
type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

// Parse extracts structured data from plain text resume using deterministic rules
func (p *Parser) Parse(rawText string) *ParseResult {
	normalized := normalizeText(rawText)
	lines := strings.Split(rawText, "\n")

	return &ParseResult{
		Skills:         p.extractSkills(normalized),
		Companies:      p.extractCompanies(normalized, lines),
		Projects:       p.extractProjects(lines),
		Education:      p.extractEducation(lines),
		Experience:     p.extractExperience(lines),
		Certifications: p.extractCertifications(normalized),
	}
}

func (p *Parser) extractSkills(text string) []string {
	found := make(map[string]bool)
	var skills []string

	for _, keywords := range TechSkillDictionary {
		for _, keyword := range keywords {
			// Use word boundary matching
			pattern := `(?i)\b` + regexp.QuoteMeta(keyword) + `\b`
			matched, _ := regexp.MatchString(pattern, text)
			if matched && !found[strings.ToLower(keyword)] {
				canonical := canonicalizeSkill(keyword)
				if !found[canonical] {
					found[canonical] = true
					skills = append(skills, canonical)
				}
			}
		}
	}

	sort.Strings(skills)
	return skills
}

func (p *Parser) extractCompanies(text string, lines []string) []string {
	found := make(map[string]bool)
	var companies []string

	// Match known companies
	for _, company := range knownCompanies {
		pattern := `(?i)\b` + regexp.QuoteMeta(company) + `\b`
		matched, _ := regexp.MatchString(pattern, text)
		if matched && !found[strings.ToLower(company)] {
			found[strings.ToLower(company)] = true
			companies = append(companies, titleCase(company))
		}
	}

	// Look for "at Company" or "@ Company" patterns in experience sections
	atPattern := regexp.MustCompile(`(?i)(?:at|@)\s+([A-Z][a-zA-Z\s&,\.]+?)(?:\s*[-,|]|\s*\(|\n|$)`)
	for _, line := range lines {
		matches := atPattern.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			if len(m) > 1 {
				company := strings.TrimSpace(m[1])
				if len(company) > 2 && len(company) < 80 && !found[strings.ToLower(company)] {
					found[strings.ToLower(company)] = true
					companies = append(companies, company)
				}
			}
		}
	}

	return unique(companies)
}

func (p *Parser) extractProjects(lines []string) []string {
	var projects []string
	inProjectSection := false
	projectHeaderRe := regexp.MustCompile(`(?i)^(projects?|personal projects?|side projects?|open source)\s*:?\s*$`)
	projectLineRe := regexp.MustCompile(`(?i)^\s*[•\-\*\d\.]+\s+([A-Za-z][\w\s\-\.]+)`)
	nextSectionRe := regexp.MustCompile(`(?i)^(experience|education|skills|certifications?|awards?|publications?)\s*:?\s*$`)

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if projectHeaderRe.MatchString(trimmed) {
			inProjectSection = true
			continue
		}
		if inProjectSection && nextSectionRe.MatchString(trimmed) {
			inProjectSection = false
			continue
		}
		if inProjectSection {
			matches := projectLineRe.FindStringSubmatch(trimmed)
			if len(matches) > 1 {
				project := strings.TrimSpace(matches[1])
				if len(project) > 3 {
					projects = append(projects, project)
				}
			}
		}
	}
	return projects
}

func (p *Parser) extractEducation(lines []string) []map[string]interface{} {
	var education []map[string]interface{}
	inEduSection := false
	eduHeaderRe := regexp.MustCompile(`(?i)^(education|academic|qualification)\s*:?\s*$`)
	nextSectionRe := regexp.MustCompile(`(?i)^(experience|skills|projects?|certifications?|work|employment)\s*:?\s*$`)
	yearRe := regexp.MustCompile(`\b(19|20)\d{2}\b`)
	degreeRe := regexp.MustCompile(`(?i)(b\.?tech|b\.?e\.?|m\.?tech|m\.?e\.?|b\.?sc|m\.?sc|mba|phd|ph\.d|bca|mca|bachelor|master|doctorate)`)

	var currentEdu map[string]interface{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			if currentEdu != nil {
				education = append(education, currentEdu)
				currentEdu = nil
			}
			continue
		}
		if eduHeaderRe.MatchString(trimmed) {
			inEduSection = true
			continue
		}
		if inEduSection && nextSectionRe.MatchString(trimmed) {
			if currentEdu != nil {
				education = append(education, currentEdu)
				currentEdu = nil
			}
			inEduSection = false
			continue
		}
		if inEduSection {
			if degreeRe.MatchString(trimmed) {
				if currentEdu != nil {
					education = append(education, currentEdu)
				}
				currentEdu = map[string]interface{}{
					"degree":      trimmed,
					"institution": "",
					"year":        "",
				}
				years := yearRe.FindAllString(trimmed, -1)
				if len(years) > 0 {
					currentEdu["year"] = years[len(years)-1]
				}
			} else if currentEdu != nil && currentEdu["institution"] == "" {
				currentEdu["institution"] = trimmed
			}
		}
	}
	if currentEdu != nil {
		education = append(education, currentEdu)
	}
	return education
}

func (p *Parser) extractExperience(lines []string) []map[string]interface{} {
	var experience []map[string]interface{}
	inExpSection := false
	expHeaderRe := regexp.MustCompile(`(?i)^(experience|work experience|professional experience|employment|work history)\s*:?\s*$`)
	nextSectionRe := regexp.MustCompile(`(?i)^(education|skills|projects?|certifications?|awards?)\s*:?\s*$`)
	yearRangeRe := regexp.MustCompile(`((?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\.?\s+)?(19|20)\d{2}\s*[-–—to]+\s*(present|current|((?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*\.?\s+)?(19|20)\d{2})`)
	roleRe := regexp.MustCompile(`(?i)^(senior|lead|principal|staff|junior|associate|mid-level|sr\.?|jr\.?)?\s*(software|backend|frontend|full.?stack|devops|cloud|data|ml|platform|site reliability|mobile|android|ios|web)\s*(engineer|developer|architect|manager|lead|analyst|scientist|consultant)`)

	var currentExp map[string]interface{}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if expHeaderRe.MatchString(trimmed) {
			inExpSection = true
			continue
		}
		if inExpSection && nextSectionRe.MatchString(trimmed) {
			if currentExp != nil {
				experience = append(experience, currentExp)
				currentExp = nil
			}
			inExpSection = false
			continue
		}
		if inExpSection {
			if yearRangeRe.MatchString(trimmed) || roleRe.MatchString(trimmed) {
				if currentExp != nil {
					experience = append(experience, currentExp)
				}
				currentExp = map[string]interface{}{
					"title":    trimmed,
					"company":  "",
					"duration": "",
				}
				yearMatch := yearRangeRe.FindString(trimmed)
				if yearMatch != "" {
					currentExp["duration"] = yearMatch
				}
			} else if currentExp != nil && currentExp["company"] == "" && len(trimmed) > 2 {
				currentExp["company"] = trimmed
			}
		}
	}
	if currentExp != nil {
		experience = append(experience, currentExp)
	}
	return experience
}

func (p *Parser) extractCertifications(text string) []string {
	var certs []string
	found := make(map[string]bool)

	for _, pattern := range certificationPatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(text, -1)
		for _, m := range matches {
			key := strings.ToLower(m)
			if !found[key] {
				found[key] = true
				certs = append(certs, strings.TrimSpace(m))
			}
		}
	}
	return certs
}

// Helper functions

func normalizeText(text string) string {
	// Replace special chars with spaces for matching
	var b strings.Builder
	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '+' || r == '#' || r == ' ' || r == '\n' || r == '\t' {
			b.WriteRune(r)
		} else {
			b.WriteRune(' ')
		}
	}
	return b.String()
}

func canonicalizeSkill(skill string) string {
	skill = strings.TrimSpace(skill)
	// Normalize well-known aliases
	aliases := map[string]string{
		"reactjs":    "React",
		"react.js":   "React",
		"react":      "React",
		"vuejs":      "Vue.js",
		"vue.js":     "Vue.js",
		"vue":        "Vue.js",
		"nodejs":     "Node.js",
		"node.js":    "Node.js",
		"springboot": "Spring Boot",
		"go":         "Go",
		"golang":     "Go",
		"k8s":        "Kubernetes",
		"postgres":   "PostgreSQL",
		"postgresql": "PostgreSQL",
		"js":         "JavaScript",
		"ts":         "TypeScript",
		"py":         "Python",
	}
	lower := strings.ToLower(skill)
	if canonical, ok := aliases[lower]; ok {
		return canonical
	}
	return titleCase(skill)
}

func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func unique(ss []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range ss {
		key := strings.ToLower(s)
		if !seen[key] {
			seen[key] = true
			result = append(result, s)
		}
	}
	return result
}
