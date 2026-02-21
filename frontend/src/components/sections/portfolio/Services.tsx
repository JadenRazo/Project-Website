import { motion } from 'framer-motion'
import { Globe, Server, Code, Database } from 'lucide-react'

const services = [
  {
    icon: Globe,
    title: 'Web Development',
    description: 'React, Next.js, Astro, TypeScript, and modern frameworks.',
  },
  {
    icon: Server,
    title: 'Host & Server Management',
    description: 'Linux servers, Nginx, Docker, CI/CD, and cloud infrastructure.',
  },
  {
    icon: Code,
    title: 'Backend Development',
    description: 'Python, Go, Node.js, and Rust for scalable APIs and services.',
  },
  {
    icon: Database,
    title: 'Database Management',
    description: 'PostgreSQL, MongoDB, Redis, and database architecture.',
  },
]

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.1,
      delayChildren: 0.2,
    },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    y: 0,
    transition: { duration: 0.5, ease: [0.22, 1, 0.36, 1] },
  },
}

export default function Services() {
  return (
    <section id="services" className="relative w-full py-12 sm:py-16 md:py-20">
      <div className="w-full flex flex-col items-center px-4 sm:px-6 lg:px-8">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          transition={{ duration: 0.6 }}
          className="mb-6 sm:mb-8 text-center"
        >
          <h2 className="text-2xl sm:text-3xl md:text-4xl font-bold mb-2 text-text-primary">
            What I <span className="gradient-text">Do</span>
          </h2>
          <p className="text-text-secondary text-sm sm:text-base">
            End-to-end development services to bring your vision to life
          </p>
        </motion.div>

        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true }}
          className="grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 w-full max-w-5xl"
        >
          {services.map((service) => (
            <motion.div
              key={service.title}
              variants={itemVariants}
              className="group glass-card p-3 sm:p-4 md:p-6 text-center hover:border-primary/50"
            >
              <div className="w-10 h-10 sm:w-12 sm:h-12 rounded-xl bg-primary/10 flex items-center justify-center mx-auto mb-3 group-hover:bg-primary/20 transition-colors">
                <service.icon className="w-5 h-5 sm:w-6 sm:h-6 text-primary" />
              </div>
              <h3 className="text-xs sm:text-sm md:text-base font-bold text-text-primary mb-1">{service.title}</h3>
              <p className="text-text-secondary text-[10px] sm:text-xs leading-relaxed hidden sm:block">
                {service.description}
              </p>
            </motion.div>
          ))}
        </motion.div>

        <motion.p
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.6, duration: 0.5 }}
          className="text-text-muted text-[10px] sm:text-xs mt-6 sm:mt-8 text-center"
        >
          Plus performance optimization, security audits, and ongoing support
        </motion.p>
      </div>
    </section>
  )
}
